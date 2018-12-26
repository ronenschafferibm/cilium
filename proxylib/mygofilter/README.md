1. Start the dev VM and ssh to it
    ```
    contrib/vagrant/start.sh
    vagrant ssh
    ```
    
1. Build and install my-go-filter
    ```
    pushd .
    cd proxylib
    make
    sudo make install
    popd
    ```

1. Start cilium-agent
    ```
    sudo systemctl start cilium
    ```

1. Run cilium-envoy independently
    ```
    sudo cilium-envoy --base-id 1 -c ./cilium-envoy-my.yaml --component-log-level filter:trace
    ``` 
    
1. From a new terminal, run pythonServer
    ```
    python pythonServer.py
    ```

1. From a new terminal, send requests to cilium-envoy
    ```
    netcat 127.0.0.1 1234 < content.txt
    ```
    
1. Check cilium-envoy logs:
    ```
    
    ```