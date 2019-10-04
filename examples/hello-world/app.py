import os
import flask

app = flask.Flask(__name__)


@app.route('/')
def hello():
    return 'Hello World!\n' + os.popen('nvidia-smi').read()


if __name__ == '__main__':
    os.environ['FLASK_ENV'] = 'development'
    port = int(os.environ.get('PORT', 8080))
    app.run(debug=True, host='0.0.0.0', port=port)
