Working directory: {{ run "pwd" }}.
Contains the following files: {{ run "ls" "-la" | nindent 4 }}
