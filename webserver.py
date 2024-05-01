from flask import Flask, send_from_directory, send_file

app = Flask("main")

host_frontend = True
if(host_frontend):
    @app.route('/')
    def send_index_html():
        return send_file('frontend/build/index.html')

    @app.route('/<path:path>')
    def send_frontend(path):
        return send_from_directory('frontend/build', path)

@app.route('/action_pics/<path:path>')
def send_action_pics(path):
    return send_from_directory('aktionen', path)

@app.route('/extern_pics/<path:path>')
def send_extern_pics(path):
    return send_from_directory('externContent', path)