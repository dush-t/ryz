import os
import sys

def run_command(command):
    print(command)
    return os.WEXITSTATUS(os.system(command))

def get_program_name(program_file):
    return os.path.basename(program_file).rstrip('.p4')

def run_compile_bmv2(program_file):
    compiler_args = []
    compiler_args.append('--p4v 16')

    # Compile the program.
    output_file = get_program_name(program_file) + '.json'
    compiler_args.append('"%s"' % program_file)
    compiler_args.append('-o "%s"' % output_file)
    rv = run_command('p4c-bm2-ss %s' % ' '.join(compiler_args))

    if rv != 0:
        print 'Compile failed.'
        sys.exit(1)

    print 'Compilation successful.\n'
    return output_file