#!/usr/bin/env bash
temporal workflow delete --query="ExecutionStatus != 'Running' and WorkflowType='EventWorkflow'" --reason "removing old" -y
