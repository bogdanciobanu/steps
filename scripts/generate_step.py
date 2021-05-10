import os
import sys

import yaml
from cookiecutter.main import cookiecutter


def kebab_to_camel(word):
    return ''.join(x.capitalize() or '-' for x in word.split('-'))


def snake_to_camel(word):
    return ''.join(x.capitalize() or '_' for x in word.split('_'))


def envs_to_go_structs(envs):
    arg_line = ""

    for e in envs:
        is_required = ""

        if e["required"]:
            is_required = ",required"
        arg_line += snake_to_camel(e["name"])
        arg_line += " {0} ".format(e["type"])
        arg_line += "`env:\"{0}{1}\"`".format(e["name"], is_required)
        arg_line += "\n"

    return arg_line


def outputs_to_go_struct(outputs):
    output_line = ""

    for o in outputs:
        if o["name"] == "api_object":
            continue

        output_line += snake_to_camel(o["name"])
        output_line += " {0} ".format(o["type"])
        output_line += "`json:\"{0}\"`".format(o["name"])
        output_line += "\n"

    return output_line


if __name__ == '__main__':
    manifest_path = sys.argv[1]
    path_list = manifest_path.split(os.sep)
    step_path, manifest_name = os.path.split(manifest_path)
    step_relpath = os.path.relpath(step_path, 'steps')
    kebab_normalized_step_name = step_path.replace(os.sep, "-")
    normalized_step_name = kebab_to_camel(path_list[-3] + "-" + path_list[-2])

    with open(manifest_path, 'r') as stream:
        y = yaml.safe_load(stream)

        step_args = envs_to_go_structs(y["envs"])
        step_output = outputs_to_go_struct(y["outputs"])

        extra_context = {'step_args': step_args, 'step_output': step_output,
                         'step_path': step_path, 'step_relpath': step_relpath, 'kebab_normalized_step_name': kebab_normalized_step_name,
                         'normalized_step_name': normalized_step_name}

        print("Generating step using the following configurations: " + str(extra_context))

        cookiecutter('scripts/step-template/',
                     extra_context=extra_context,
                     overwrite_if_exists=True,
                     no_input=True)
