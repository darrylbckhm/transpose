import re
import boto3
import scrapy
import logging
from scrapy import Selector
from boto3.dynamodb.types import TypeSerializer
from botocore.exceptions import ClientError
from lib.constants import RE_URLS
from crawler.spiders.real_estate.items import RealEstateItem


logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.DEBUG)


URLS=['https://www.re.cr/en/costa-rica-real-estate-for-sale/search-properties?form.search-properties.widgets.listing_type:list=rs&form.search-properties.widgets.listing_type:list=cs&form.search-properties.widgets.listing_type:list=ll&b_start:int=0&form.search-properties.buttons.search=']

class RealEstateSpider(scrapy.Spider):
    name = "RealEstate"

    def start_requests(self, urls=URLS):
        for url in urls:
            yield scrapy.Request(url=url, callback=self.parse)

    def parse(self, response):
        _meta = response.xpath('//figure').extract()
        _attrs = response.css('.tileItem').xpath("//div[contains(@class, 'listing__')]").extract()

        _properties = [[_] for _ in _meta]

        i = j = 0
        for _ in range(len(_attrs)):
            if i > 0 and '<div class="listing__price">' in _attrs[i]:
                j += 1
            if j >= len(_properties) or i >= len(_attrs):
                break
            _properties[j].append(_attrs[i])
            i += 1
                
        for prop in _properties:
            _meta = [_.split('\"') for _ in prop[0].split('=')]
            _url = ['URL', _meta[1][1].strip()]
            _title = ['Title', _meta[2][1].strip()]
            _img = ['Image', _meta[3][1].strip()]

            cleanr = re.compile('<.*?>')
            cleandata = [_url, _title, _img]
            for i in range(1, len(prop[1:])):
                _prop = [s.strip() for s in re.sub(cleanr, '', prop[i]).split("\n") if s.strip()]
                if len(_prop) == 1 and _prop[0] == 'Property Type':
                    _prop.append('Other')
                cleandata.append(_prop)

            p = dict(cleandata)
            p = {k.replace(' ', '_').lower(): v for k,v in p.items()}
            yield p

        nextpage = response.css('span.next a::attr(href)').get()
        if nextpage:
            yield scrapy.Request(url=nextpage, callback=self.parse)
