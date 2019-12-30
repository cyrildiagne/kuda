import os
import flask

app = flask.Flask(__name__)

# The only route which returns the result of nvidia-smi.
@app.route('/')
def hello():
    return 'Hello GPU!\n\n' + os.popen('nvidia-smi').read()


# Development.
if __name__ == '__main__':
    port = int(os.environ.get('PORT', 80))
    app.run(debug=True, host='0.0.0.0', port=port)
