import React, { Component, PropTypes } from 'react';

export default class SmoothieGraph extends Component {
  addRandomValueToDataSet(time, dataSet) {
      dataSet.append(time, Math.random());
  }
  initHost(hostId) {

    // Initialize an empty TimeSeries for each CPU.
    var cpuDataSet = new TimeSeries();
    var seriesOption =  { strokeStyle: 'rgba(255, 0, 0, 1)', fillStyle: 'rgba(255, 0, 0, 0.1)', lineWidth: 3 };

    var now = new Date().getTime();
    for (var t = now - 1000 * 50; t <= now; t += 1000) {
      this.addRandomValueToDataSet(t, cpuDataSet);
    }
    // Every second, simulate a new set of readings being taken from each CPU.
    setInterval(function() {
      this.addRandomValueToDataSet(new Date().getTime(), cpuDataSet);
    }.bind(this), 1000);

    // Build the timeline
    var timeline = new SmoothieChart({ millisPerPixel: 20, grid: { strokeStyle: '#555555', lineWidth: 1, millisPerLine: 1000, verticalSections: 4 }});
    timeline.addTimeSeries(cpuDataSet, seriesOption);
    timeline.streamTo(React.findDOMNode(this.refs[hostId + 'Cpu']), 1000);
  }
  componentDidMount() {
    this.initHost('host1');
  }
  render() {
    var canvasStyle = {
      width: '100%'
    };
    return (
      <canvas id="host1Cpu" ref="host1Cpu" height="100" style={canvasStyle}></canvas>
    );
  }
}

SmoothieGraph.propTypes = {
  id: PropTypes.string.isRequired
}
