import requests
import json
import sys

# Here we define our query as a multi-line string
format = str(sys.argv[2])
if format == "NONE":
  format = ""
else:
  format = ", format: " + str(sys.argv[2])

query = '''
query ($page: Int, $perPage: Int, $search: String, $id: Int) {
    Page (page: $page, perPage: $perPage) {
        pageInfo {
            lastPage
        }
        media (search: $search, id: $id''' + format + ''') {
            id
            status
            format
            episodes
            duration
            chapters
            volumes
            averageScore
            source
            startDate {
              day
              month
              year
            }
            endDate {
              day
              month
              year
            }
            genres
            
            studios(isMain : true) {
              nodes {
                name
              }
            }

            title {
                romaji
                english
            }

            description
            coverImage {
              extraLarge
            }
            siteUrl
            externalLinks {
              site
              url
            }
            relations {
              edges {
                relationType
              }
              nodes {
                title {
                  english
                  romaji
                }
                id
              }
            }
            nextAiringEpisode {
              timeUntilAiring
              episode
            }

        }
    }
}
'''

variables = {
    'search': str(sys.argv[1]),
    'page': 1,
    'perPage': 10
}

id = str(sys.argv[3])
if id != "NONE":
  print(int(id))
  variables = {
    'page': 1,
    'perPage': 10,
    'id': int(id)
  }

print(variables)
url = 'https://graphql.anilist.co'

response = requests.post(url, json={'query': query, 'variables': variables})

with open('config/response.json', 'w', encoding='utf-8') as f:
    json.dump(response.json(), f, ensure_ascii=False, indent=4)