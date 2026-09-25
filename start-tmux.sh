#!/bin/bash

# Ensure the 'recordalerting' session exists
if ! tmux has-session -t recordalerting 2>/dev/null; then
  # Create a new session named 'recordalerting', detached
  cd /workspaces/recordalerting
  tmux new-session -d -s recordalerting
  
  # Split the window horizontally (-h)
  # The left pane will remain a terminal
  # The right pane will run 'gh dash'
  tmux split-window -h -t recordalerting "gh dash"
fi
