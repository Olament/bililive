import React from 'react';
import 'antd/dist/antd.css';
import './index.css';
import './App.css';
// import { List } from 'antd';
import Broadcast from "./Broadcast";
import VirtualScroller from "virtual-scroller/react";

export default class CardList extends React.Component {
    state = {
        loading: true,
        count: 0,
        data: [],
    };

    componentDidMount() {
        clearInterval(this.timer)
        this.getData(res => {
            this.setState({
                count: res.count,
                data: res.list,
                loading: false,
            })
        })
        this.timer = setInterval(() => {
            this.getData(res => {
                this.setState({
                    count: res.count,
                    data: res.list,
                    loading: false,
                })
            })
        }, 30000)
    };

    componentWillUnmount() {
        this.timer && clearInterval(this.timer)
    }

    getData = callback => {
        fetch('https://livevup.com/api/online')
            .then(res => res.json())
            .then(res => callback(res))
    }

    render() {
        return (
            <VirtualScroller
                id="broadcasts"
                items={this.state.data}
                itemComponent={({children}) => (
                    <div style={{padding: '14px 14px 8px'}}>
                        <Broadcast item={children}></Broadcast>
                    </div>
                )}
                getColumnsCount={(container) => (Math.min(Math.floor(container.getWidth() / 240), 4))}
            />
        );
    }
}