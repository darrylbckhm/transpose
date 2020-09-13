import boto3

class AWSClientBase:
    def __init__(self, service, endpoint=None, aws_access_key_id=None, aws_secret_access_key=None):
        self.service = service
        self.endpoint = endpoint
        self.aws_access_key_id = aws_access_key_id
        self.aws_secret_access_key = aws_secret_access_key

    def get_client(self):
        return boto3.client(self.service, endpoint_url=self.endpoint)
