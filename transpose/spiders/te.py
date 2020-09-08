


import re
import boto3
import scrapy
import logging
from scrapy import Selector
from boto3.dynamodb.types import TypeSerializer
from botocore.exceptions import ClientError
from marshmallow_dataclass import dataclass as marshmallow_dataclass

logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.DEBUG)

@marshmallow_dataclass
class Property:
    title: str
    url: str
    image: str
    price: str
    listing_type: str
    property_type: str
    status: str
    location: str
    location_type: str
    view_type: str

class teSpider(scrapy.Spider):
    name = "te"

    def start_requests(self):
        urls = [
            'https://tradingeconomics.com/costa-rica/indicators'
        ]
        for url in urls:
            yield scrapy.Request(url=url, callback=self.parse)

    def parse(self, response):
        return response
