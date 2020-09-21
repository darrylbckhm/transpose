#!/usr/local/bin/python3

import boto3
from lib.transpose import Transposer
from twisted.internet import reactor
from scrapy.crawler import CrawlerProcess
from crawler.spiders.Transpose import TransposeSpider
from scrapy.utils.project import get_project_settings



def main():
    crawler = CrawlerProcess(get_project_settings())
    crawler.crawl(TransposeSpider)
    crawler.start()
    #reactor.run() # the script will block here

    #transposer = Transposer()
    #transposer.transpose()

if __name__ == "__main__":
    main()
