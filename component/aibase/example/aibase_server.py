# Copyright (C) 2024 wwhai
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
from concurrent import futures
import signal
import grpc
import aibase_pb2
import aibase_pb2_grpc


class AIBaseServiceServicer(aibase_pb2_grpc.AIBaseServiceServicer):
    def Call(self, request, context):
        # 处理普通请求
        result = b"OK"
        return aibase_pb2.CallResponse(result=result)

    def Stream(self, request_iterator, context):
        # 处理流式请求
        for request in request_iterator:
            result = b"OK"
            yield aibase_pb2.StreamResponse(result=result)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    aibase_pb2_grpc.add_AIBaseServiceServicer_to_server(AIBaseServiceServicer(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()


def handle_sigint(signum, frame):
    raise SystemExit(1)


def handle_sigterm(signum, frame):
    raise SystemExit(1)


if __name__ == "__main__":
    signal.signal(signal.SIGINT, handle_sigint)
    signal.signal(signal.SIGTERM, handle_sigterm)
    serve()
