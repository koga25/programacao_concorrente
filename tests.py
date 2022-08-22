import subprocess as subp
import matplotlib.pyplot as plt
import numpy as np

concurrent_clients_tests = [1, 5, 10, 20, 40, 80]

# TCP Server
subp.check_call(["go", "run", "./tcp/tcpServer.go", "1234"])

# UDP Server
subp.check_call(["go", "run", "./udp/udpServer.go", "1234"])

tcp_rtt = []
udp_rtt = []

for client_num in concurrent_clients_tests:
    print(f"Now testing for {client_num} concurrent clients.")
    # TCP Clients
    for i in range(client_num-1):
       subp.check_call(["go", "run", "./tcp/tcpClient.go", "localhost:1234", "false"])
    # Last client needs to measure time RTT and return it as a mean of durations.
    current_client_num_rtt = subp.check_call(["go", "run", "./tcp/tcpClient.go", "localhost:1234", "true"])
    tcp_rtt.append(current_client_num_rtt)

    # UDP Clients
    for i in range(client_num-1):
       subp.check_call(["go", "run", "./udp/udpClient.go", "localhost:1234", "false"])
    # Last client needs to measure time RTT and return it as a mean of durations.
    current_client_num_rtt = subp.check_call(["go", "run", "./udp/udpClient.go", "localhost:1234", "true"])
    udp_rtt.append(current_client_num_rtt)


# x = np.arange(len(concurrent_clients_tests))  # the label locations
# width = 0.35  # the width of the bars
  
# tcp_bar = plt.bar(x - width/2, tcp_rtt, width, label = 'TCP')
# udp_bar = plt.bar(x + width/2, udp_rtt, width, label = 'UDP')

# plt.xticks(x, concurrent_clients_tests)
# plt.xlabel("Quantidade de clientes concorrentes")
# plt.ylabel("Média de RTT entre 10000 requisições")
# plt.title("Diferenças de performance entre aplicações usando TCP e UDP em diversas quantidades de clientes concorrentes")
# plt.legend()

# plt.bar_label(tcp_bar, padding=3)
# plt.bar_label(udp_bar, padding=3)

# plt.show()

        

    