# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: aibase.proto
# Protobuf Python Version: 4.25.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0c\x61ibase.proto\x12\x06\x61ibase\"\x1b\n\x0b\x43\x61llRequest\x12\x0c\n\x04\x64\x61ta\x18\x01 \x01(\x0c\"\x1e\n\x0c\x43\x61llResponse\x12\x0e\n\x06result\x18\x01 \x01(\x0c\"\x1d\n\rStreamRequest\x12\x0c\n\x04\x64\x61ta\x18\x01 \x01(\x0c\" \n\x0eStreamResponse\x12\x0e\n\x06result\x18\x01 \x01(\x0c\x32\x83\x01\n\rAIBaseService\x12\x33\n\x04\x43\x61ll\x12\x13.aibase.CallRequest\x1a\x14.aibase.CallResponse\"\x00\x12=\n\x06Stream\x12\x15.aibase.StreamRequest\x1a\x16.aibase.StreamResponse\"\x00(\x01\x30\x01\x42\x1d\n\x06\x61ibaseB\x06\x61ibaseP\x00Z\t./;aibaseb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'aibase_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  _globals['DESCRIPTOR']._options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\006aibaseB\006aibaseP\000Z\t./;aibase'
  _globals['_CALLREQUEST']._serialized_start=24
  _globals['_CALLREQUEST']._serialized_end=51
  _globals['_CALLRESPONSE']._serialized_start=53
  _globals['_CALLRESPONSE']._serialized_end=83
  _globals['_STREAMREQUEST']._serialized_start=85
  _globals['_STREAMREQUEST']._serialized_end=114
  _globals['_STREAMRESPONSE']._serialized_start=116
  _globals['_STREAMRESPONSE']._serialized_end=148
  _globals['_AIBASESERVICE']._serialized_start=151
  _globals['_AIBASESERVICE']._serialized_end=282
# @@protoc_insertion_point(module_scope)
