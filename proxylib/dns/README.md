This is filter demonstrates how to use cilium feature to write envoy filters in Go, in order write a dns filter that parses dns requests (over TCP currently because envoy doesn't support UDP).

# My Environment
I have a repo of cilium (https://github.com/cilium/cilium.git) on my Windows host machine under:
`C:\Users\ronensch\cilium`

that directory is shared with a Xubuntu guest machine under: `~/go/src/cilium`

I couldn't do all work inside the guest machine since VirtualBox can't run inside a virtual machine.


# Build & Run

1. Change `daemon/Makefile` links to copied version since windows filesystem doesn't support symlinks
    Comment out:
    ```
    $(foreach link,$(LINKS), ln -f -s $(TARGET) $(link);)
    ```
	and add:
	```
	$(foreach link,$(LINKS), cp $(TARGET) $(link);)
	```

1. Start the virtual machine builder from host machine:
    ```
    ./contrib/vagrant/start.sh
    ```

1. Log in to the virtual machine builder
    ```
    vagrant ssh
    ```

1. From inside the virtual machine builder, build and install the filters
    ```
    ./rebuild.bash
    ```
    
1. From inside the virtual machine builder, run cilium-envoy with config for dns filter
    ```
    sudo cilium-envoy --base-id 1 -c ./cilium-envoy-dns.yaml --component-log-level filter:trace
    ```

1. Open a new terminal and ssh to the virtual machine builder
    ```
    vagrant ssh
    ```
    
1. From the new sessions, issue a dns request over tcp to the running cilium-envoy
    ```
    dig +tcp @localhost -p 53  abc.com
    ```
    
    The output should look like:
    ```
    ;; Connection to ::1#53(::1) for abc.com failed: connection refused.
    
    ; <<>> DiG 9.11.3-1ubuntu1.1-Ubuntu <<>> +tcp @localhost -p 53 abc.com
    ; (2 servers found)
    ;; global options: +cmd
    ;; Got answer:
    ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 6451
    ;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
    
    ;; OPT PSEUDOSECTION:
    ; EDNS: version: 0, flags:; udp: 512
    ;; QUESTION SECTION:
    ;abc.com.                       IN      A
    
    ;; ANSWER SECTION:
    abc.com.                146     IN      A       199.181.132.250
    
    ;; Query time: 135 msec
    ;; SERVER: 127.0.0.1#53(127.0.0.1)
    ;; WHEN: Mon Jan 14 12:19:49 UTC 2019
    ;; MSG SIZE  rcvd: 52
    ```
    
    Envoy's log (in the first terminal) should look like:
    ```
    [2019-01-14 12:21:02.972][007038][info][main] [external/envoy/source/server/server.cc:203] initializing epoch 0 (hot restart version=10.200.16384.127.options=capacity=16384, num_slots=8209 hash=228984379728933363 size=2654312)
    [2019-01-14 12:21:02.972][007038][info][main] [external/envoy/source/server/server.cc:205] statically linked extensions:
    [2019-01-14 12:21:02.972][007038][info][main] [external/envoy/source/server/server.cc:207]   access_loggers: envoy.file_access_log,envoy.http_grpc_access_log
    [2019-01-14 12:21:02.972][007038][info][main] [external/envoy/source/server/server.cc:210]   filters.http: cilium.l7policy,envoy.buffer,envoy.cors,envoy.ext_authz,envoy.fault,envoy.filters.http.header_to_metadata,envoy.filters.http.jwt_authn,envoy.filters.http.rbac,envoy.grpc_http1_bridge,envoy.grpc_json_transcoder,envoy.grpc_web,envoy.gzip,envoy.health_check,envoy.http_dynamo_filter,envoy.ip_tagging,envoy.lua,envoy.rate_limit,envoy.router,envoy.squash,istio_authn,jwt-auth,mixer
    [2019-01-14 12:21:02.972][007038][info][main] [external/envoy/source/server/server.cc:213]   filters.listener: cilium.bpf_metadata,envoy.listener.original_dst,envoy.listener.proxy_protocol,envoy.listener.tls_inspector
    [2019-01-14 12:21:02.972][007038][info][main] [external/envoy/source/server/server.cc:216]   filters.network: cilium.network,envoy.client_ssl_auth,envoy.echo,envoy.ext_authz,envoy.filters.network.dubbo_proxy,envoy.filters.network.rbac,envoy.filters.network.sni_cluster,envoy.filters.network.tcp_cluster_rewrite,envoy.filters.network.thrift_proxy,envoy.http_connection_manager,envoy.mongo_proxy,envoy.ratelimit,envoy.redis_proxy,envoy.tcp_proxy,mixer
    [2019-01-14 12:21:02.972][007038][info][main] [external/envoy/source/server/server.cc:218]   stat_sinks: envoy.dog_statsd,envoy.metrics_service,envoy.stat_sinks.hystrix,envoy.statsd
    [2019-01-14 12:21:02.972][007038][info][main] [external/envoy/source/server/server.cc:220]   tracers: envoy.dynamic.ot,envoy.lightstep,envoy.tracers.datadog,envoy.zipkin
    [2019-01-14 12:21:02.972][007038][info][main] [external/envoy/source/server/server.cc:223]   transport_sockets.downstream: envoy.transport_sockets.alts,envoy.transport_sockets.capture,raw_buffer,tls
    [2019-01-14 12:21:02.972][007038][info][main] [external/envoy/source/server/server.cc:226]   transport_sockets.upstream: envoy.transport_sockets.alts,envoy.transport_sockets.capture,raw_buffer,tls
    [2019-01-14 12:21:02.979][007038][info][main] [external/envoy/source/server/server.cc:268] admin address: 127.0.0.1:8081
    [2019-01-14 12:21:02.980][007038][info][config] [external/envoy/source/server/configuration_impl.cc:50] loading 0 static secret(s)
    [2019-01-14 12:21:02.980][007038][info][config] [external/envoy/source/server/configuration_impl.cc:56] loading 2 cluster(s)
    [2019-01-14 12:21:02.982][007038][info][upstream] [external/envoy/source/common/upstream/cluster_manager_impl.cc:136] cm init: all clusters initialized
    [2019-01-14 12:21:02.982][007038][info][config] [external/envoy/source/server/configuration_impl.cc:61] loading 2 listener(s)
    [2019-01-14 12:21:02.983][007038][info][filter] [cilium/proxylib.cc:17] GoFilter: Opening go module libcilium.so
    INFO[0000] init(): Registering cassandraParserFactory
    INFO[0000] NPDS: Registering L7 rule parser: cassandra
    INFO[0000] init(): Registering dnsParserFactory
    INFO[0000] init(): Registering memcacheParserFactory
    INFO[0000] NPDS: Registering L7 rule parser: memcache
    INFO[0000] init(): Registering r2d2ParserFactory
    INFO[0000] NPDS: Registering L7 rule parser: r2d2
    INFO[0000] init(): Registering blockParserFactory
    INFO[0000] init(): Registering headerParserFactory
    INFO[0000] NPDS: Registering L7 rule parser: test.headerparser
    INFO[0000] init(): Registering lineParserFactory
    INFO[0000] init(): Registering PasserParserFactory
    INFO[0000] Opened new library instance 1
    [2019-01-14 12:21:02.996][007038][trace][filter] [./cilium/socket_option.h:38] Set socket (23) option SO_MARK to b00 (magic mark: b00, id: 0, cluster: 0)
    [2019-01-14 12:21:02.996][007038][trace][filter] [./cilium/socket_option.h:19] Skipping setting socket (23) option SO_MARK, state != STATE_PREBIND
    [2019-01-14 12:21:02.996][007038][info][config] [external/envoy/source/server/configuration_impl.cc:94] loading tracing configuration
    [2019-01-14 12:21:02.996][007038][info][config] [external/envoy/source/server/configuration_impl.cc:112] loading stats sink configuration
    [2019-01-14 12:21:02.996][007038][info][main] [external/envoy/source/server/server.cc:428] all clusters initialized. initializing init manager
    [2019-01-14 12:21:02.996][007038][info][config] [external/envoy/source/server/listener_manager_impl.cc:908] all dependencies initialized. starting workers
    [2019-01-14 12:21:02.996][007038][info][main] [external/envoy/source/server/server.cc:456] starting main dispatch loop
    [2019-01-14 12:21:02.999][007050][trace][filter] [./cilium/socket_option.h:19] Skipping setting socket (23) option SO_MARK, state != STATE_PREBIND
    [2019-01-14 12:21:02.999][007049][trace][filter] [./cilium/socket_option.h:19] Skipping setting socket (23) option SO_MARK, state != STATE_PREBIND
    [2019-01-14 12:21:10.641][007049][debug][filter] [./cilium/socket_option.h:57] Cilium SocketOption(): source_identity: 2, destination_identity: 2, ingress: false, port: 53, proxy_port: 0
    [2019-01-14 12:21:10.641][007049][trace][filter] [cilium/bpf_metadata.cc:132] cilium.bpf_metadata (egress): GOT metadata for new connection
    [2019-01-14 12:21:10.641][007049][debug][filter] [external/envoy/source/common/tcp_proxy/tcp_proxy.cc:193] [C1] new tcp proxy session
    [2019-01-14 12:21:10.641][007049][debug][filter] [cilium/network_filter.cc:72] Cilium Network: onNewConnection
    [2019-01-14 12:21:10.641][007049][debug][filter] [cilium/network_filter.cc:89] [C1] Cilium Network: No proxymap
    [2019-01-14 12:21:10.641][007049][trace][filter] [cilium/proxylib.cc:111] [C1] GoFilter: Calling go module
    DEBU[0007] dnsParserFactory: Create: &{0xc0000d49b0 1 false 2 2 127.0.0.1:58035 127.0.0.1:53 127.0.0.1 53 dns <nil> 0x55555741c048 0x55555741c4a0}
    [2019-01-14 12:21:10.642][007049][debug][filter] [external/envoy/source/common/tcp_proxy/tcp_proxy.cc:335] [C1] Creating connection to cluster dns-google
    [2019-01-14 12:21:10.642][007049][trace][filter] [./cilium/socket_option.h:38] Set socket (26) option SO_MARK to 20b00 (magic mark: b00, id: 2, cluster: 0)
    [2019-01-14 12:21:10.706][007049][debug][filter] [external/envoy/source/common/tcp_proxy/tcp_proxy.cc:521] TCP:onUpstreamEvent(), requestedServerName:
    [2019-01-14 12:21:10.707][007049][trace][filter] [cilium/proxylib.cc:224] [C1] Cilium Network::OnIO: Calling go module with 50 bytes of data
    DEBU[0007] *-*-*-OnData-dnsMsg: [158 31 1 32 0 1 0 0 0 0 0 1 3 97 98 99 3 99 111 109 0 0 1 0 1 0 0 41 16 0 0 0 0 0 0 12 0 10 0 8 201 223 190 45 161 110 161 153]
    DEBU[0007] *-*-*-OnData-data: [0 48 158 31 1 32 0 1 0 0 0 0 0 1 3 97 98 99 3 99 111 109 0 0 1 0 1 0 0 41 16 0 0 0 0 0 0 12 0 10 0 8 201 223 190 45 161 110 161 153]
    DEBU[0007] *-*-*-OnData-m: ;; opcode: QUERY, status: NOERROR, id: 40479
    ;; flags: rd ad; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 1
    
    ;; QUESTION SECTION:
    ;abc.com.       IN       A
    
    ;; ADDITIONAL SECTION:
    
    ;; OPT PSEUDOSECTION:
    ; EDNS: version 0; flags: ; udp: 4096
    ; COOKIE: c9dfbe2da16ea199
    
    DEBU[0007] *-*-*-OnData-m.Question: [{abc.com. 1 1}]
    DEBU[0007] *-*-*-========= *NOT* GOOGLE.COM
    
    DEBU[0007] *-*-*-OnData-len(data): 50
    DEBU[0007] *-*-*-OnData-reply: false
    [2019-01-14 12:21:10.707][007049][trace][filter] [cilium/proxylib.cc:226] [C1] Cilium Network::OnIO: 'go_on_data' returned 0, ops(1)
    [2019-01-14 12:21:10.707][007049][debug][filter] [cilium/proxylib.cc:251] [C1] Cilium Network::OnIO: FILTEROP_PASS: 50 bytes
    [2019-01-14 12:21:10.707][007049][debug][filter] [cilium/proxylib.cc:305] [C1] Cilium Network::OnIO: Output on return: [2019-01-14 12:21:10.780][007049][trace][filter] [external/envoy/source/common/tcp_proxy/tcp_proxy.cc:480] [C1] upstream connection received 54 bytes, end_stream=false
    [2019-01-14 12:21:10.780][007049][trace][filter] [cilium/proxylib.cc:224] [C1] Cilium Network::OnIO: Calling go module with 54 bytes of data
    DEBU[0007] *-*-*-OnData-dnsMsg: [158 31 129 128 0 1 0 1 0 0 0 1 3 97 98 99 3 99 111 109 0 0 1 0 1 192 12 0 1 0 1 0 0 0 64 0 4 199 181 132 250 0 0 41 2 0 0 0 0 0 0 0]
    DEBU[0007] *-*-*-OnData-data: [0 52 158 31 129 128 0 1 0 1 0 0 0 1 3 97 98 99 3 99 111 109 0 0 1 0 1 192 12 0 1 0 1 0 0 0 64 0 4 199 181 132 250 0 0 41 2 0 0 0 0 0 0 0]
    DEBU[0007] *-*-*-OnData-m: ;; opcode: QUERY, status: NOERROR, id: 40479
    ;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
    
    ;; QUESTION SECTION:
    ;abc.com.       IN       A
    
    ;; ANSWER SECTION:
    abc.com.        64      IN      A       199.181.132.250
    
    ;; ADDITIONAL SECTION:
    
    ;; OPT PSEUDOSECTION:
    ; EDNS: version 0; flags: ; udp: 512
    
    DEBU[0007] *-*-*-OnData-m.Question: [{abc.com. 1 1}]
    DEBU[0007] *-*-*-========= *NOT* GOOGLE.COM
    
    DEBU[0007] *-*-*-OnData-m.Answer: net.IP{0xc7, 0xb5, 0x84, 0xfa}
    DEBU[0007] *-*-*-OnData-m.Answer: &dns.A{Hdr:dns.RR_Header{Name:"abc.com.", Rrtype:0x1, Class:0x1, Ttl:0x40, Rdlength:0x4}, A:net.IP{0xc7, 0xb5, 0x84, 0xfa}}
    DEBU[0007] *-*-*-OnData-len(data): 54
    DEBU[0007] *-*-*-OnData-reply: true
    [2019-01-14 12:21:10.781][007049][trace][filter] [cilium/proxylib.cc:226] [C1] Cilium Network::OnIO: 'go_on_data' returned 0, ops(1)
    [2019-01-14 12:21:10.781][007049][debug][filter] [cilium/proxylib.cc:251] [C1] Cilium Network::OnIO: FILTEROP_PASS: 54 bytes
    [2019-01-14 12:21:10.781][007049][debug][filter] [cilium/proxylib.cc:305] [C1] Cilium Network::OnIO: Output on return: [2019-01-14 12:21:10.847][007049][trace][filter] [external/envoy/source/common/tcp_proxy/tcp_proxy.cc:480] [C1] upstream connection received 0 bytes, end_stream=true
    [2019-01-14 12:21:10.847][007049][trace][filter] [cilium/proxylib.cc:224] [C1] Cilium Network::OnIO: Calling go module with 0 bytes of data
    [2019-01-14 12:21:10.847][007049][trace][filter] [cilium/proxylib.cc:226] [C1] Cilium Network::OnIO: 'go_on_data' returned 0, ops(0)
    [2019-01-14 12:21:10.847][007049][debug][filter] [cilium/proxylib.cc:305] [C1] Cilium Network::OnIO: Output on return:
    [2019-01-14 12:21:10.847][007049][trace][filter] [cilium/network_filter.cc:147] [C1] Cilium Network::OnWrite: 'GoFilter::OnIO' returned 0
    ```

### See also:
- [Envoy Go Extensions](https://cilium.readthedocs.io/en/v1.3/envoy/extensions/)
- [Developer / Contributor Guide](https://cilium.readthedocs.io/en/v1.3/contributing/)
- https://github.com/ronenschafferibm/proxy-1/blob/local-dev/README.local_dev.md

## Further points
- Consider creating a docker image instead of vagrant virtual machine in order to build envoy's go filters.
- Extend `network.cilium` filter to enable sharing data between go filters and non go filters ([dynamicMetadata](https://www.envoyproxy.io/docs/envoy/latest/configuration/well_known_dynamic_metadata) mechanism).
 We need this capability so we capability so dns-parser go filter could pass data to istio's mixer filter, that will log and enforce policies.
- Discuss how we should integrate go filters in istio-proxy repo