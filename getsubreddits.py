import requests


def main():
  num_of_pages = 40
  for i in range(num_of_pages):
    r = requests.get("http://redditlist.com/all?page=%s" % i)
    text = r.text.split("\n")
    
    for x in text:
      #if True:
      if "subreddit-url" in x:
         print(x)

if __name__ == "__main__":
  main()
