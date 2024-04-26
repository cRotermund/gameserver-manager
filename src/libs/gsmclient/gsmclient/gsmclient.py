import requests
import time
from requests_auth_aws_sigv4 import AWSSigV4

ENDPOINT_URL = "https://vd5vhqweprchm3fixxqb4emphe0zpolz.lambda-url.us-east-1.on.aws/"
SERVER_WAIT_TIMEOUT = 30

class GSMClient:
    def __init__(self, key: str, secret: str, region: str):
        self.key = key
        self.secret = secret
        self.region = region
        self.service = "lambda"

    #PRIVATES prefixed, "--"
    ######################################################################
    ######################################################################
    ######################################################################

    def __get_auth(self):
        return AWSSigV4(
            aws_access_key_id = self.key,
            aws_secret_access_key = self.secret,
            region = self.region,
            service = self.service
        )

    def __send(
        self,
        body: object
    ):
        sig = self.__get_auth()

        r = requests.request(
            "POST",
            ENDPOINT_URL,
            json = body,
            auth = sig
        )
        return r

    # END PRIVATES prefixed, "--"
    ######################################################################
    ######################################################################
    ######################################################################

    # PUBLICS
    ######################################################################
    ######################################################################
    ######################################################################

    def start(self):
        body = {"action": "start"}
        r = self.__send(body)

    def stop(self):
        body = {"action": "stop"}
        r = self.__send(body)

    def status(self):
        body = { "action" : "status" }
        r = self.__send(body)
        return r.json()

    def wait_for_status(self, desired: str, onpoll: callable):
        reached = False
        started_at = time.time()
        while not reached:
            s = self.status()["status"]
            if onpoll is not None:
                onpoll(s)

            reached = (s == desired)
            elapsed = (time.time() - started_at)

            if not reached and elapsed > SERVER_WAIT_TIMEOUT:
                raise Exception("Timed out waiting for server status.")

            if not reached:
                time.sleep(1)

    # END PUBLICS
    ######################################################################
    ######################################################################
    ######################################################################