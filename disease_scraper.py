import requests
from string import ascii_lowercase
from bs4 import BeautifulSoup
import random
import time
import csv

browser_agents = [
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:73.0) Gecko/20100101 Firefox/73.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:53.0) Gecko/20100101 Firefox/53.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36"]
url = "https://www.nhp.gov.in/disease-a-z/{letter_placeholder}"
all_diseases = []
for letter in ascii_lowercase:
    user_agent = browser_agents[random.randint(0, len(browser_agents) - 1)]
    rq = requests.get(url.format(letter_placeholder=letter), verify=False, timeout=5,headers={'User-Agent':user_agent})
    soup = BeautifulSoup(rq.content, "html.parser")
    tags = soup.select('.all-disease a li')
    diseases = [tag.text for tag in tags]
    all_diseases.extend(diseases)
    time.sleep(random.randint(100, 400) / 1000)
with open('data/diseases.csv', "w", newline="\n") as f:
    for disease in all_diseases:
        f.write(disease)