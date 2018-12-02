import os
from subprocess import run, PIPE

def markdown_to_html(request):
    f = request.files.to_dict().get('doc').read().decode()

    p = run(['./pandoc', '--from=markdown', '--to=html5'], stdout=PIPE, input=f, encoding='utf-8')

    return (p.stdout, 200, {'Content-Type': 'text/html; charset=UTF-8'})
