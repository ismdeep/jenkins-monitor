import copy
import io
import sys
import json
import logging
import os.path
import time

import requests
from clint.textui import colored
from dateutil import parser as dateutil_parser
from prettytable import PrettyTable
from pytz import timezone


def args_exists(__key__):
    if len(sys.argv) <= 1:
        return False
    for item in sys.argv[1:]:
        if __key__ == item:
            return True
    return False


def get_argv(__key__):
    for i in range(1, len(sys.argv) - 1):
        if __key__ == sys.argv[i]:
            return sys.argv[i + 1]
    return ""


def help_msg():
    return '''Usage: python3 doraemon-shell.py [--help]|[-c config.json [show <pods|apps>]|login|status|[reload pods]]
    e.g.
        --help                                       Show help information
        -c config.json login                         Login
        -c config.json status                        Show status
        -c config.json show pods                     Show running pods information
        -c config.json show apps                     Show apps information
        -c config.json reload pods                   Reload pods
        -c config.json reload pods --wait            Reload pods wait to quit
    '''


def beauty_time(time_raw):
    created_at = dateutil_parser.parse(time_raw)
    return created_at.astimezone(timezone('Asia/Shanghai'))


def time_elapse_text(time_raw):
    created_at = dateutil_parser.parse(time_raw)
    s = int(time.time()) - created_at.timestamp()
    if s < 60:
        return f'{s} seconds'
    minutes = int(s / 60)
    if minutes < 60:
        return f'{minutes} minutes'
    hours = int(minutes / 60)
    if hours < 24:
        return f'{hours} hours'
    days = int(hours / 24)
    return f'{days} days'


class K8sClient(object):
    site = ''
    token = ''
    jweToken = ''
    config = None
    service_name = ''
    apps = []
    data_path = None

    def __init__(self):
        self.token = ''
        logging.captureWarnings(True)

    def load_config(self, config_path):
        self.config = json.load(open(config_path))
        self.site = self.config['site']
        self.token = self.config['token']
        self.service_name = self.config['service_name']
        self.apps = self.config['apps']
        self.data_path = self.config['data']
        if self.data_path is None or self.data_path == '':
            print(colored.red("[ERROR] load config"))
            exit(0)
        if not os.path.exists(self.data_path):
            os.mkdir(self.data_path)
        if os.path.isfile("{}/jweToken".format(self.data_path)):
            with open("{}/jweToken".format(self.data_path), 'r') as f:
                try:
                    self.jweToken = f.readline()
                except io.UnsupportedOperation:
                    self.jweToken = ''

    def login(self):
        req = requests.get(
            url='{}/api/v1/csrftoken/login'.format(self.site),
            verify=False,
        )
        data = json.loads(req.text)
        if 'token' not in data:
            return False
        csrf_token = data['token']
        req = requests.post(
            url='{}/api/v1/login'.format(self.site),
            data=json.dumps({
                'token': self.token,
            }),
            headers={
                'x-csrf-token': csrf_token,
                'content-type': 'application/json',
            },
            verify=False,
        )
        data = json.loads(req.text)
        if len(data['errors']) <= 0:
            self.jweToken = data['jweToken']
            with open("{}/jweToken".format(self.data_path), 'w') as f:
                f.write(self.jweToken)
            return True
        else:
            return False

    def get_pods(self):
        req = requests.get(
            url='{}/api/v1/pod/{}?itemsPerPage=10&page=1&sortBy=d,creationTimestamp'.format(self.site,
                                                                                            self.service_name),
            headers={
                'accept': 'application/json, text/plain, */*',
                'jwetoken': '{}'.format(self.jweToken),
            },
            verify=False,
        )
        if req.text.strip().find("MSG_LOGIN_UNAUTHORIZED_ERROR") >= 0:
            print(colored.red("[ERROR] MSG_LOGIN_UNAUTHORIZED_ERROR"))
            return None
        data = json.loads(req.text)
        if len(data['errors']) >= 1:
            return None
        return data['pods']

    def show_pods(self):
        pods = self.get_pods()
        if pods is None:
            return
        x = PrettyTable()
        x.field_names = ["Pod", 'App Name', "Status", "CreatedAt", "Ago"]
        for pod in pods:
            status = pod['podStatus']['status']
            if status == "Running":
                status = colored.green(status).color_str
            x.add_row(
                [pod['objectMeta']['name'], pod['objectMeta']['labels']['app'], status,
                 beauty_time(pod['objectMeta']['creationTimestamp']),
                 time_elapse_text(pod['objectMeta']['creationTimestamp'])])
        print(x)

    def get_apps(self):
        req = requests.get(
            url='{}/api/v1/deployment/doraemon?itemsPerPage=10&page=1&sortBy=d,creationTimestamp'.format(self.site),
            headers={
                'accept': 'application/json, text/plain, */*',
                'jwetoken': '{}'.format(self.jweToken),
            },
            verify=False,
        )
        if req.text.strip().find("MSG_LOGIN_UNAUTHORIZED_ERROR") >= 0:
            print(colored.red("[ERROR] MSG_LOGIN_UNAUTHORIZED_ERROR"))
            return None
        data = json.loads(req.text)
        if len(data['errors']) >= 1:
            return None
        return data

    def show_apps(self):
        data = self.get_apps()
        if data is None:
            return
        x = PrettyTable()
        x.field_names = ["Deployment Name", "App Name", "Image", "CreatedAt", "Ago"]
        for deployment in data['deployments']:
            app_name = ''
            if 'kubectl.kubernetes.io/last-applied-configuration' in deployment['objectMeta']['annotations']:
                last_applied_config_str = deployment['objectMeta']['annotations'][
                    'kubectl.kubernetes.io/last-applied-configuration']
                obj = json.loads(last_applied_config_str)
                if app_name == '':
                    try:
                        app_name = obj['spec']['template']['metadata']['labels']['app']
                    except KeyError:
                        pass
                if app_name == '':
                    try:
                        app_name = obj['spec']['selector']['matchLabels']['app']
                    except KeyError:
                        pass

            x.add_row([deployment['objectMeta']['name'], app_name, ','.join(deployment['containerImages']),
                       beauty_time(deployment['objectMeta']['creationTimestamp']),
                       time_elapse_text(deployment['objectMeta']['creationTimestamp'])])
        print(x)

    def is_running(self, __pod_name__):
        pods = self.get_pods()
        for pod in pods:
            if pod['objectMeta']['name'] == __pod_name__:
                return True
        return False

    def reload_pods(self, wait=False):
        pods = self.get_pods()
        reload_pods = []
        results = {}

        for pod in pods:
            if pod['objectMeta']['labels']['app'] in self.apps:
                reload_pods.append(pod['objectMeta']['name'])

        for reload_pod_name in reload_pods:
            req = requests.delete(
                url='{}/api/v1/_raw/pod/namespace/doraemon/name/{}'.format(self.site, reload_pod_name),
                headers={
                    'accept': 'application/json, text/plain, */*',
                    'jwetoken': '{}'.format(self.jweToken),
                },
                verify=False,
            )
            results[reload_pod_name] = req.text.strip() == ""

        for reload_pod_name in reload_pods:
            if not results[reload_pod_name]:
                print(colored.red('Reload [{}] Failed.'.format(reload_pod_name)))
            else:
                if wait:
                    while self.is_running(reload_pod_name):
                        time.sleep(1)
                print(colored.green('Reload [{}] Success.'.format(reload_pod_name)))

        return True

    def show_status(self):
        loaded_config = copy.deepcopy(self.config)
        loaded_config['token'] = '******'
        print(json.dumps(loaded_config, indent=4))

        print()

        if self.jweToken is None or self.jweToken == '':
            print(colored.red('Is Not Login'))
        else:
            print(colored.green('Login Success'))


def main():
    client = K8sClient()
    client.load_config(get_argv("-c"))
    client.login()
    client.reload_pods(wait=True)


if __name__ == '__main__':
    main()
