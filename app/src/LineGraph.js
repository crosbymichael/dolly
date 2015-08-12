import React, { Component } from 'react';

var lineChartData = [
  // The first layer
  {
    label: "Layer 1",
    values: [ {time: 1370044800, y: 100}, {time: 1370044801, y: 1000}, {time: 1370044800, y: 100}, {time: 1370044801, y: 1000} ]
  },
  // The second layer
  {
    label: "Layer 2",
    values: [ {time: 1370044800, y: 78}, {time: 1370044801, y: 98}, {time: 1370044800, y: 100}, {time: 1370044801, y: 1000} ]
  }
];

export default class LineGraph extends Component {
  componentDidMount() {
    $('#lineChart').epoch({
      type: 'time.line',
      data: lineChartData
    });
  }
  render() {
    var chartStyle = {
      width: '100%',
      height: '100px',
      borderRadius: '3px',
      border: '1px solid black'
    };
    return (
      <div id="lineChart" style={chartStyle}></div>
    );
  }
}
