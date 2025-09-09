clean:
	rm -rf go dart kotlin python/wampprotofbs/gen/*

gen: clean gen-py gen-go gen-dart gen-kotlin

gen-py:
	flatc -o python/wampprotofbs/gen --python flatbuffers/*.fbs --gen-all
	mv python/wampprotofbs/gen/wampproto/* python/wampprotofbs/gen/
	rm -rf python/wampprotofbs/gen/wampproto

gen-go:
	flatc -o go --go flatbuffers/message.fbs --gen-all

gen-dart:
	flatc -o dart --dart flatbuffers/message.fbs --gen-all

gen-kotlin:
	flatc -o kotlin --kotlin flatbuffers/message.fbs --gen-all

verify:
	./.venv/bin/python create.py
