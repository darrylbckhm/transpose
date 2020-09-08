import os
import sys
import setuptools


if not hasattr(sys, 'real_prefix'):
    sys.path.append(os.path.dirname(__file__))

with open("README.md", "r") as fh:
    long_description = fh.read()

setuptools.setup(
    name="transpose",
    version="0.0.1",
    description="transpose",
    long_description=long_description,
    long_description_content_type="text/markdown",
    packages=setuptools.find_packages(),
    include_package_date=True,
    install_requires=[
        'boto3',
        'pytest',
        'scrapy',
        'marshmallow_dataclass',
    ],
    entry_points={
          'console_scripts': ['ftf=bin.transpose:main'],
      },
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires='>=3.6',
)
