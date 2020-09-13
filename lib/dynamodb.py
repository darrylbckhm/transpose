import logging
from boto3.dynamodb.types import TypeSerializer
from botocore.exceptions import ClientError
from lib.aws import AWSClientBase


class DynamoDBClient(AWSClientBase):
    def __init__(self, client=None, endpoint=None):
        self.service = 'dynamodb'
        self.endpoint = endpoint
        self.client = client or self.get_client()

    def put_item_v2(self, table_name, item, dataclass=None):
        serializer = TypeSerializer()
        try:
            response = self.client.put_item(
                TableName=table_name,
                Item={k: serializer.serialize(v) for k, v in item.items() if v != ""}
            )
        except ClientError as err:
            raise err
        else:
            return response
