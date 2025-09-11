install_uv:
	@if ! command -v uv >/dev/null 2>&1; then \
  		curl -LsSf https://astral.sh/uv/install.sh | sh; \
  	fi

setup:
	make install_uv
	make gen
	uv venv
	uv pip install .[test,publish]

lint:
	.venv/bin/ruff format .

check-lint:
	.venv/bin/ruff check .

test:
	.venv/bin/pytest -v tests/

build-wheel:
	pip wheel --no-deps -w dist .

clean:
	rm -rf .venv build dist wampprotofbs.egg-info

publish-build:
	rm -rf ./dist ./build ./wampprotofbs.egg-info/
	.venv/bin/python -m build --sdist
	.venv/bin/twine check dist/*
	@echo ========================================================
	@echo
	@echo now run .venv/bin/twine upload dist/newly_created.tar.gz

gen:
	$(MAKE) -C .. gen-py

install-flatc:
	sudo apt install cmake -y
	sudo apt install flatbuffers-compiler -y
