if [[ -z $(git status --porcelain) ]]; then
            echo 0
            exit 0
          fi
          echo 1