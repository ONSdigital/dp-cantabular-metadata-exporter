Feature: Cantabular-Metadata-Exporter

  Background:

    And the following metadata document with dataset id "cantabular-example-1", edition "2021" and version "1" is available from dp-dataset-api:
      """
      {
        "dimensions": [
          {
            "label": "City",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/city",
            "id": "city",
            "name": "City"
          },
          {
            "label": "Number of siblings (3 mappings)",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/siblings",
            "id": "siblings",
            "name": "Number of siblings (3 mappings)"
          },
          {
            "label": "Sex",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/sex",
            "id": "sex",
            "name": "Sex"
          }
        ],
        "distribution": [
          "json",
          "csvw",
          "txt"
        ],
        "downloads": {},
        "release_date": "2021-11-19T00:00:00.000Z",
        "title": "Test Cantabular Dataset Published",
        "headers": [
          "cantabular_table",
          "city",
          "siblings_3",
          "sex"
        ]
      }
      """

    And the following metadata document with dataset id "cantabular-example-2", edition "2021" and version "1" is available from dp-dataset-api:
      """
      {
        "dimensions": [
          {
            "label": "City",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/city",
            "id": "city",
            "name": "City"
          },
          {
            "label": "Number of siblings (3 mappings)",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/siblings",
            "id": "siblings",
            "name": "Number of siblings (3 mappings)"
          },
          {
            "label": "Sex",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/sex",
            "id": "sex",
            "name": "Sex"
          }
        ],
        "distribution": [
          "json",
          "csvw",
          "txt"
        ],
        "downloads": {},
        "release_date": "2021-11-19T00:00:00.000Z",
        "title": "Test Cantabular Dataset Unpublished",
        "headers": [
          "cantabular_table",
          "city",
          "siblings_3",
          "sex"
        ]
      }
      """

    And the following version document with dataset id "cantabular-example-1", edition "2021" and version "1" is available from dp-dataset-api:
      """
      {
        "alerts": [],
        "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
        "dimensions": [
          {
            "label": "City",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/city",
            "id": "city",
            "name": "City"
          },
          {
            "label": "Number of siblings (3 mappings)", 
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/siblings",
            "id": "siblings",
            "name": "Number of siblings (3 mappings)"
          },
          {
            "label": "Sex",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/sex",
            "id": "sex",
            "name": "Sex"
          }
        ],
        "edition": "2021",
        "id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
        "links": {
          "dataset": {
            "href": "http://dp-dataset-api:22000/datasets/cantabular-example-1",
            "id": "cantabular-example-1"
          },
          "dimensions": {},
          "edition": {
            "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021",
            "id": "2021"
          },
          "self": {
            "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1"
          }
        },
        "release_date": "2021-11-19T00:00:00.000Z",
        "state": "published",
        "usage_notes": [],
        "version": 1
      }
      """

    And the following version document with dataset id "cantabular-example-2", edition "2021" and version "1" is available from dp-dataset-api:
      """
      {
        "alerts": [],
        "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
        "dimensions": [
          {
            "label": "City",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/city",
            "id": "city",
            "name": "City"
          },
          {
            "label": "Number of siblings (3 mappings)",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/siblings",
            "id": "siblings",
            "name": "Number of siblings (3 mappings)"
          },
          {
            "label": "Sex",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/sex",
            "id": "sex",
            "name": "Sex"
          }
        ],
        "edition": "2021",
        "id": "c733977d-a2ca-4596-9cb1-08a6e724858a",
        "links": {
          "dataset": {
            "href": "http://dp-dataset-api:22000/datasets/cantabular-example-1",
            "id": "cantabular-example-1"
          },
          "dimensions": {},
          "edition": {
            "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021",
            "id": "2021"
          },
          "self": {
            "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1"
          }
        },
        "release_date": "2021-11-19T00:00:00.000Z",
        "state": "edition-confirmed",
        "usage_notes": [],
        "version": 1
      }
      """

    And the following dimensions are available from dataset "cantabular-example-1" edition "2021" version "1":
      """
      {
        "items": [
          {
            "links": {
              "code_list": {
              "href": "http://localhost:22400/code-lists/city",
              "id": "city"
            },
            "options": {
              "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1/dimensions/city/options",
              "id": "city"
            },
            "version": {
              "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1"
            }
          },
          "name": "city"
          },
          {
            "links": {
              "code_list": {
                "href": "http://localhost:22400/code-lists/sex",
                "id": "sex"
              },
              "options": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1/dimensions/sex/options",
                "id": "sex"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1"
              }
            },
            "name": "sex"
          },
          {
            "links": {
              "code_list": {
                "href": "http://localhost:22400/code-lists/siblings_3",
                "id": "siblings_3"
              },
              "options": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1/dimensions/siblings_3/options",
                "id": "siblings_3"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1"
              }
            },
            "name": "siblings_3"
          }
        ],
        "count": 3,
        "offset": 0,
        "limit": 20,
        "total_count": 3
      }
      """

    And the following dimensions are available from dataset "cantabular-example-2" edition "2021" version "1":
      """
      {
        "items": [
          {
            "links": {
              "code_list": {
              "href": "http://localhost:22400/code-lists/city",
              "id": "city"
            },
            "options": {
              "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1/dimensions/city/options",
              "id": "city"
            },
            "version": {
              "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1"
            }
          },
          "name": "city"
          },
          {
            "links": {
              "code_list": {
                "href": "http://localhost:22400/code-lists/sex",
                "id": "sex"
              },
              "options": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1/dimensions/sex/options",
                "id": "sex"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1"
              }
            },
            "name": "sex"
          },
          {
            "links": {
              "code_list": {
                "href": "http://localhost:22400/code-lists/siblings_3",
                "id": "siblings_3"
              },
              "options": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1/dimensions/siblings_3/options",
                "id": "siblings_3"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1"
              }
            },
            "name": "siblings_3"
          }
        ],
        "count": 3,
        "offset": 0,
        "limit": 20,
        "total_count": 3
      }
      """

    And the following options response is available for dimension "city" for dataset "cantabular-example-1" edition "2021" version "1" with query params "offset=0&limit=1000":
      """
      {
        "items": [
          {
            "label": "London",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/city/codes/0",
                "id": "0"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/city",
                "id": "city"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "city",
            "option": "0"
          },
          {
            "label": "Liverpool",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/city/codes/1",
                "id": "1"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/city",
                "id": "city"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "city",
            "option": "1"
          },
          {
            "label": "Belfast",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/city/codes/2",
                "id": "2"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/city",
                "id": "city"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "city",
            "option": "2"
          }
        ],
        "count": 3,
        "offset": 0,
        "limit": 1000,
        "total_count": 3
      }
      """

    And the following options response is available for dimension "sex" for dataset "cantabular-example-1" edition "2021" version "1" with query params "offset=0&limit=1000":
      """
      {
        "items": [
          {
            "label": "Male",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/sex/codes/0",
                "id": "0"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/sex",
                "id": "sex"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "sex",
            "option": "0"
          },
          {
            "label": "Female",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/sex/codes/1",
                "id": "1"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/sex",
                "id": "sex"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "sex",
            "option": "1"
          }
        ],
        "count": 2,
        "offset": 0,
        "limit": 1000,
        "total_count": 2
      }
      """
    And the following options response is available for dimension "siblings_3" for dataset "cantabular-example-1" edition "2021" version "1" with query params "offset=0&limit=1000":
    """
    {
      "items": [
        {
          "label": "No siblings",
          "links": {
            "code": {
              "href": "http://localhost:22400/code-lists/siblings_3/codes/0",
              "id": "0"
            },
            "code_list": {
              "href": "http://localhost:22400/code-lists/siblings_3",
              "id": "siblings_3"
            },
            "version": {
              "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
              "id": "1"
            }
          },
          "dimension": "siblings_3",
          "option": "0"
        },
        {
          "label": "1 or 2 siblings",
          "links": {
            "code": {
              "href": "http://localhost:22400/code-lists/siblings_3/codes/1-2",
              "id": "1-2"
            },
            "code_list": {
              "href": "http://localhost:22400/code-lists/siblings_3",
              "id": "siblings_3"
            },
            "version": {
              "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
              "id": "1"
            }
          },
          "dimension": "siblings_3",
          "option": "1-2"
        },
        {
          "label": "3 or more siblings",
          "links": {
            "code": {
              "href": "http://localhost:22400/code-lists/siblings_3/codes/3+",
              "id": "3+"
            },
            "code_list": {
              "href": "http://localhost:22400/code-lists/siblings_3",
              "id": "siblings_3"
            },
            "version": {
              "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
              "id": "1"
            }
          },
          "dimension": "siblings_3",
          "option": "3+"
        }
      ],
      "count": 3,
      "offset": 0,
      "limit": 1000,
      "total_count": 3
    }
    """
    And the following options response is available for dimension "city" for dataset "cantabular-example-2" edition "2021" version "1" with query params "offset=0&limit=1000":
      """
      {
        "items": [
          {
            "label": "London",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/city/codes/0",
                "id": "0"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/city",
                "id": "city"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-2/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "city",
            "option": "0"
          },
          {
            "label": "Liverpool",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/city/codes/1",
                "id": "1"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/city",
                "id": "city"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-2/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "city",
            "option": "1"
          },
          {
            "label": "Belfast",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/city/codes/2",
                "id": "2"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/city",
                "id": "city"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-2/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "city",
            "option": "2"
          }
        ],
        "count": 3,
        "offset": 0,
        "limit": 1000,
        "total_count": 3
      }
      """
    And the following options response is available for dimension "sex" for dataset "cantabular-example-2" edition "2021" version "1" with query params "offset=0&limit=1000":
      """
      {
        "items": [
          {
            "label": "Male",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/sex/codes/0",
                "id": "0"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/sex",
                "id": "sex"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-2/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "sex",
            "option": "0"
          },
          {
            "label": "Female",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/sex/codes/1",
                "id": "1"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/sex",
                "id": "sex"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-2/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "sex",
            "option": "1"
          }
        ],
        "count": 2,
        "offset": 0,
        "limit": 1000,
        "total_count": 2
      }
      """

    And the following options response is available for dimension "siblings_3" for dataset "cantabular-example-2" edition "2021" version "1" with query params "offset=0&limit=1000":
      """
      {
        "items": [
          {
            "label": "No siblings",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/siblings_3/codes/0",
                "id": "0"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/siblings_3",
                "id": "siblings_3"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-2/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "siblings_3",
            "option": "0"
          },
          {
            "label": "1 or 2 siblings",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/siblings_3/codes/1-2",
                "id": "1-2"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/siblings_3",
                "id": "siblings_3"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-2/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "siblings_3",
            "option": "1-2"
          },
          {
            "label": "3 or more siblings",
            "links": {
              "code": {
                "href": "http://localhost:22400/code-lists/siblings_3/codes/3+",
                "id": "3+"
              },
              "code_list": {
                "href": "http://localhost:22400/code-lists/siblings_3",
                "id": "siblings_3"
              },
              "version": {
                "href": "http://localhost:22000/datasets/cantabular-example-2/editions/2021/versions/1",
                "id": "1"
              }
            },
            "dimension": "siblings_3",
            "option": "3+"
          }
        ],
        "count": 3,
        "offset": 0,
        "limit": 1000,
        "total_count": 3
      }
      """

    Given the following version with dataset id "cantabular-example-2", edition "2021" and version "1" will be updated to dp-dataset-api:
      """
      {
        "alerts": null,
        "collection_id": "",
        "downloads": {
          "CSVW": {
            "href": "http://localhost:23600/downloads/datasets/cantabular-example-2/editions/2021/versions/1.csv-metadata.json",
            "size": "643",
            "private": "http://minio:9000/dp-cantabular-metadata-exporter-priv/datasets/cantabular-example-2-2021-1.csvw"
          },
          "TXT": {
            "href": "http://localhost:23600/downloads/datasets/cantabular-example-2/editions/2021/versions/1.txt",
            "size": "503",
            "private": "http://minio:9000/dp-cantabular-metadata-exporter-priv/datasets/cantabular-example-2-2021-1.txt"
          }
        },
        "edition": "",
        "dimensions": null,
        "id": "",
        "instance_id": "",
        "latest_changes": null,
        "links": {
          "access_rights": {
            "href": ""
          },
          "dataset": {
            "href": ""
          },
          "dimensions": {
            "href": ""
          },
          "edition": {
            "href": ""
          },
          "editions": {
            "href": ""
          },
          "latest_version": {
            "href": ""
          },
          "versions": {
            "href": ""
          },
          "self": {
            "href": ""
          },
          "code_list": {
            "href": ""
          },
          "options": {
            "href": ""
          },
          "version": {
            "href": ""
          },
          "code": {
            "href": ""
          },
          "taxonomy": {
            "href": ""
          },
          "job": {
            "href": ""
          }
        },
        "release_date": "",
        "state": "",
        "temporal": null,
        "version": 0
      }
      """

    And the following version with dataset id "cantabular-example-1", edition "2021" and version "1" will be updated to dp-dataset-api:
      """
      {
        "alerts": null,
        "collection_id": "",
        "downloads": {
          "CSVW": {
            "href": "http://localhost:23600/downloads/datasets/cantabular-example-1/editions/2021/versions/1.csv-metadata.json",
            "size": "641",
            "public": "http://minio:9000/dp-cantabular-metadata-exporter-pub/datasets/cantabular-example-1-2021-1.csvw"
          },
          "TXT": {
            "href": "http://localhost:23600/downloads/datasets/cantabular-example-1/editions/2021/versions/1.txt",
            "size": "499",
            "public": "http://minio:9000/dp-cantabular-metadata-exporter-pub/datasets/cantabular-example-1-2021-1.txt"
          }
        },
        "edition": "",
        "dimensions": null,
        "id": "",
        "instance_id": "",
        "latest_changes": null,
        "links": {
          "access_rights": {
            "href": ""
          },
          "dataset": {
            "href": ""
          },
          "dimensions": {
            "href": ""
          },
          "edition": {
            "href": ""
          },
          "editions": {
            "href": ""
          },
          "latest_version": {
            "href": ""
          },
          "versions": {
            "href": ""
          },
          "self": {
            "href": ""
          },
          "code_list": {
            "href": ""
          },
          "options": {
            "href": ""
          },
          "version": {
            "href": ""
          },
          "code": {
            "href": ""
          },
          "taxonomy": {
            "href": ""
          },
          "job": {
            "href": ""
          }
        },
        "release_date": "",
        "state": "",
        "temporal": null,
        "version": 0
      }
      """

  Scenario: Consuming a cantabular-metadata-export event with the correct fields and the collection is published

    Given dataset-api is healthy

    And the service starts

    When this cantabular-metadata-export event is consumed:
      """
      {
        "datasetID":   "cantabular-example-1",
        "edition":     "2021",
        "version":     "1",
        "InstanceID":  "test-instance-01",
        "RowCount":    5
      }
      """

    Then a file with filename "datasets/cantabular-example-1-2021-1.csvw" can be seen in minio bucket "dp-cantabular-metadata-exporter-pub"

    And these CSVW Created events should be produced:

      | InstanceID       | DatasetID            | Edition | Version | RowCount |
      | test-instance-01 | cantabular-example-1 | 2021    | 1       | 5        |

  Scenario: Consuming a cantabular-metadata-export event with the correct fields and the collection is not published

    Given dataset-api is healthy

    And the service starts

    When this cantabular-metadata-export event is consumed:
      """
      {
        "datasetID":   "cantabular-example-2",
        "edition":     "2021",
        "version":     "1",
        "InstanceID":  "test-instance-02",
        "RowCount":    3
      }
      """

    Then a file with filename "datasets/cantabular-example-2-2021-1.csvw" can be seen in minio bucket "dp-cantabular-metadata-exporter-priv"

    And these CSVW Created events should be produced:

      | InstanceID        | DatasetID            | Edition | Version | RowCount |
      | test-instance-02  | cantabular-example-2 | 2021    | 1       | 3        |

  Scenario: Consuming a cantabular-metadata-export event with the correct fields but a downstream service is unhealthy

    Given dataset-api is unhealthy

    And the service starts

    When this cantabular-metadata-export event is consumed:
      """
      {
        "datasetID":   "cantabular-example-1",
        "edition":     "2021",
        "version":     "1",
        "InstanceID":  "test-instance-01",
        "RowCount":    5
      }
      """

    Then no CSVW Created events should be produced