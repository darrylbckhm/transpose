# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html


# useful for handling different item types with a single interface
from itemadapter import ItemAdapter
from lib.dynamodb import DynamoDBClient
from crawler.spiders.real_estate.itemss import RealEstateItem


class WriteToDBPipeline:
    def __init__(self, dynamodb=None, table_name='local.test', endpoint='http://localhost:4566'):
        self.table_name = table_name
        self.endpoint = endpoint
        self.dynamodb = dynamodb or DynamoDBClient(endpoint=endpoint)
        self.dynamodb.client.create_table( \
                            TableName='local.test', \
                            AttributeDefinitions=[{'AttributeName': 'url', 'AttributeType': 'S'}], \
                            KeySchema=[{'AttributeName': 'url', 'KeyType': 'HASH'}], \
                            ProvisionedThroughput={'ReadCapacityUnits': 1, 'WriteCapacityUnits': 1} \
                    )

    def process_item(self, item, spider):
        real_estate_item = RealEstateItem()
        # TODO actually use items properly and set key programatically
        real_estate_item['title'] = item['title']
        real_estate_item['url'] = item['url']
        real_estate_item['image'] = item['image']
        real_estate_item['price'] = item['price']
        real_estate_item['listing_type'] = item['listing_type']
        real_estate_item['property_type'] = item['property_type']
        real_estate_item['status'] = item['status']
        real_estate_item['location'] = item['location']
        real_estate_item['location_type'] = item['location_type']
        real_estate_item['view_type'] = item['view_type']
        
        self.dynamodb.put_item_v2(self.table_name, real_estate_item, dataclass=PropertyItem)

        return item
