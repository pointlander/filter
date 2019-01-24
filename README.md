# Traffic Spike Filter

In main.go a traffic spike detector is implemented and simulated. Traffic spikes are detected regardless of the base load traffic request rate. Traffic spikes are detected by feeding the traffic rate into a [low pass filter](https://en.wikipedia.org/wiki/Low-pass_filter#Discrete-time_realization) and then computing a ratio by dividing the instantaneous rate by the filtered rate. If the ratio is greater than some threshold then there is a traffic spike. The algorithm has two parameters: the number of recorded requests stored in a ring buffer and the spike threshold.

Below is a graph of the ratio: request rate divided by the low pass filtered request rate.
![rate / filtered rate](points.png?raw=true)
