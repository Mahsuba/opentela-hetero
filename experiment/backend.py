# OpenTela-Hetero: stand-in job runner for the mode experiment.
# It answers every request (and /health) with its own name, so the reply
# shows which worker ran the job. It does no real GPU work.
# Usage: python3 experiment/backend.py PORT NAME
import http.server
import json
import sys

port, name = int(sys.argv[1]), sys.argv[2]


class Handler(http.server.BaseHTTPRequestHandler):
    def reply(self):
        n = int(self.headers.get("Content-Length") or 0)
        if n:
            self.rfile.read(n)
        body = json.dumps({"backend": name}).encode()
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    do_GET = reply
    do_POST = reply

    def log_message(self, *args):
        pass


http.server.ThreadingHTTPServer(("127.0.0.1", port), Handler).serve_forever()
