# Define here the models for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

import scrapy


class RealEstateItem(scrapy.Item):
    # define the fields for your item here like:
    # name = scrapy.Field()
    title: str = scrapy.Field()
    url: str = scrapy.Field()
    image: str = scrapy.Field()
    price: str = scrapy.Field()
    listing_type: str = scrapy.Field()
    property_type: str = scrapy.Field()
    status: str = scrapy.Field()
    location: str = scrapy.Field()
    location_type: str = scrapy.Field()
    view_type: str = scrapy.Field()
