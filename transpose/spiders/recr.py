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

DDB = boto3.client('dynamodb', endpoint_url='http://localhost:4566')

def put_item_v2(table_name, item):
    serializer = TypeSerializer()
    try:
        response = DDB.put_item(
            TableName=table_name,
            Item={k: serializer.serialize(v) for k, v in Property.Schema().dump(item).items() if v != ""}
        )
    except ClientError as err:
        raise err
    else:
        return response

class recrSpider(scrapy.Spider):
    name = "recr"

    def start_requests(self):
        urls = [
            'https://www.re.cr/en/costa-rica-real-estate-for-sale/search-properties?form.search-properties.widgets.listing_type:list=rs&form.search-properties.widgets.listing_type:list=cs&form.search-properties.widgets.listing_type:list=ll&b_start:int=0&form.search-properties.buttons.search='
        ]
        for url in urls:
            yield scrapy.Request(url=url, callback=self.parse)

    def parse(self, response):
        page = int(response.url.split("&")[-2].split('=')[-1]) / 10
        filename = 'properties-%s.html' % page
        _meta = response.xpath('//figure').extract()
        _attrs = response.css('.tileItem').xpath("//div[contains(@class, 'listing__')]").extract()

        properties = []
        for i in range(len(_meta)):
            properties.append([_meta[i]])

        i = j = 0
        for _ in range(len(_attrs)):
            if i > 0 and '<div class="listing__price">' in _attrs[i]:
                j += 1
            if j >= len(properties) or i >= len(_attrs):
                break
            properties[j].append(_attrs[i])
            i += 1
                
        response = DDB.create_table(TableName='local.test', AttributeDefinitions=[{'AttributeName': 'url', 'AttributeType': 'S'}], KeySchema=[{'AttributeName': 'url', 'KeyType': 'HASH'}], ProvisionedThroughput={'ReadCapacityUnits': 1, 'WriteCapacityUnits': 1})
        with open(filename, 'w') as f:
            for prop in properties:
                prop[0] = [_.split('\"') for _ in prop[0].split('=')]
                _url = ['URL', prop[0][1][1].strip()]
                _title = ['Title', prop[0][2][1].strip()]
                _img = ['Image', prop[0][3][1].strip()]

                cleanr = re.compile('<.*?>')
                cleandata = [_url, _title, _img]
                for i in range(1, len(prop[1:])):
                    prop[i] = [s.strip() for s in re.sub(cleanr, '', prop[i]).split("\n") if s.strip()]
                    if len(prop[i]) == 1 and prop[i][0] == 'Property Type':
                        prop[i].append('Other')
                    cleandata.append(prop[i])

                p = dict(cleandata)
                p = {k.replace(' ', '_').lower(): v for k,v in p.items()}
                p = Property(
                        p['title'],
                        p['url'],
                        p['image'],
                        p['price'],
                        p['listing_type'],
                        p['property_type'],
                        p['status'],
                        p['location'],
                        p['location_type'],
                        p['view_type'],
                    )
                put_item_v2('local.test', p)
        self.log('Saved file %s' % filename)

        nextpage = response.css('span.next a::attr(href)').get()
        if nextpage:
            yield scrapy.Request(url=nextpage, callback=self.parse)
