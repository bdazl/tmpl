Working directory: {{ run "pwd" }}.
Contains the following files: {{ run "ls" "-la" | nindent 4 }}
Exit code of "true": {{ exitCode "true" }}
Exit code of "false": {{ exitCode "false" }}
