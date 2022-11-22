# Jupyter notebooks for data analysis

## Setup

- Install [Anaconda](https://www.anaconda.com/products/distribution)
- Create Anaconda environment
  ```sh
  conda env create -f ./conda.yml
  ```
- Activate Anaconda environment
  ```sh
  conda activate e2eanalysis
  ```

To add packages:

- Update [./conda.yml](./conda.yml)
- Update environment
  ```sh
  conda env update -f ./conda.yml
  ```

To use vscode, see https://code.visualstudio.com/docs/datascience/jupyter-notebooks

### Windows

On windows it is difficult to enable Anaconda on the default PowerShell terminal.

It is best to:

- Launch the "Anaconda PowerShell terminal" from the Start prompt
- Luanch Visual Studio Code from the terminal: `code`


## Refresh all notebooks

After [fetching new data](../fetch/README.md), you can refresh all notebooks by running `.\refresh.ps1` on Windows PowerShell.