protolint
=========

protolint is the pluggable linting/fixing utility for Protocol Buffer
files (proto2+proto3)

-  Runs fast because this works without compiler.
-  Easy to follow the official style guide. The rules and the style
   guide correspond to each other exactly.

   -  Fixer automatically fixes all the possible official style guide
      violations.

-  Allows to disable rules with a comment in a Protocol Buffer file.

   -  It is useful for projects which must keep API compatibility while
      enforce the style guide as much as possible.
   -  Some rules can be automatically disabled by inserting comments to
      the spotted violations.

-  Loads plugins to contain your custom lint rules.
-  Undergone testing for all rules.
-  Many integration supports.

   -  protoc plugin
   -  Editor integration
   -  GitHub Action
   -  CI Integration

   .. rubric:: Usage in python projects
      :name: usage-in-python-projects

You can use ``protolint`` as a linter within your python projects, the
wheel ``protolint-bin`` on `pypi <https://pypi.org>`__ contains the
pre-compiled binaries for various platforms. Just add the desired
version to your ``pyproject.toml`` or ``requirements.txt``.

The wheels downloaded will contain the compiled go binaries for
``protolint`` and ``protoc-gen-protolint``. Your platform must be
compatible with the supported binary platforms.

You can add the linter configuration to the ``tools.protolint`` package
in ``pyproject.toml``.

More information
----------------

You will find more information on the `projects
homepage <https://github.com/yoheimuta/protolint>`__.
