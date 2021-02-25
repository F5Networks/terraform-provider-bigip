---
name: Bug report
about: Report a defect in the product
title: ''
labels: bug, untriaged
assignees: ''

---

<!--
When filing an issue please check to see if an issue already exists that matches your's issue.

Please open a case (https://support.f5.com/csp/article/K2633) with F5 if this is a critical issue.

-->

### Environment

 * TMOS/Bigip Version:
 * Terraform Version:
 * Terraform bigip provider Version:

### Summary

A clear and concise description of what the bug is.
Please also include information about the reproducibility and the severity/impact of the issue.

### Steps To Reproduce

Steps to reproduce the behavior:

1. Provide terraform resource config which you are facing trouble along with the output of it.

2. To get to know more about the issue, provide terraform debug logs

3. To capture debug logs, export TF_LOG variable with debug ( export TF_LOG= DEBUG ) before 
  runnning terraform apply/plan

4. As3/DO json along with the resource config( for AS3/DO resource issues )


### Expected Behavior
A clear and concise description of what you expected to happen.

### Actual Behavior
A clear and concise description of what actually happens.
Please include any applicable error output.


