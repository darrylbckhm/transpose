import json
import pytest
from tests.response_json import response_json

class TestCostaRicaRealEstate:
    def test_load_local_response(self):
        for resp in response_json:
            for k,v in resp.items():
                print(f'{k}: {v}')
            print("\n\n")
        assert False
