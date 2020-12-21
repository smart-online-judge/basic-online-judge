To run all tests,
`go test -v ./...`

Excecute `python src/config/first_start.py` in the service-related container before the first start.
It will upload nltk data, gensim model required for Python algorithms to the ~/home directory.

Python required packages list can be found in `src/config/requirements.txt`.