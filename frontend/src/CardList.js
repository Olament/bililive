import React from 'react';
import 'antd/dist/antd.css';
import './index.css';
import { List } from 'antd';
import Broadcast from "./Broadcast";


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
            <List
                grid={{
                    gutter: 16,
                    xs: 1,
                    sm: 1,
                    md: 2,
                    lg: 4,
                    xl: 4,
                    xxl: 4,
                }}
                dataSource={this.state.data}
                loading={this.state.loading}
                renderItem={item => (
                    <List.Item>
                        <Broadcast
                            item={item}
                        >
                        </Broadcast>
                    </List.Item>
                )}
            />
        );
    }
}
