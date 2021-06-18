import React from 'react';
import {List, Statistic, Row, Col, Radio} from 'antd';
import 'antd/dist/antd.css';
import './App.css'
import {withRouter} from "react-router";


class Rank extends React.Component {
    state = {
        data: [],
        loading: true,
    }

    componentDidMount() {
        this.update("income")
    }

    getData = (sortBy, callback) => {
        fetch('https://livevup.com/api/rank?sortBy=' + sortBy)
            .then(res => res.json())
            .then(res => callback(res))
    }

    update = (sortBy) => (
        this.getData(sortBy, (res) => (this.setState({
            data: res,
            loading: false
        })))
    )

    handleRadio = (e) => {
        this.setState({
            loading: true,
        })
        this.update(e.target.value.toString())
    }

    render() {
        return (
            <div>
                <span>
                    排序：
                    <Radio.Group
                        defaultValue="income"
                        buttonStyle="solid"
                        onChange={this.handleRadio}
                    >
                        <Radio.Button value="income">收入</Radio.Button>
                        <Radio.Button value="viewership">平均同接</Radio.Button>
                        <Radio.Button value="paid">场均付费</Radio.Button>
                        <Radio.Button value="duration">时长</Radio.Button>
                    </Radio.Group>
                </span>
                <List
                    itemLayout="horizontal"
                    size="large"
                    loading={this.state.loading}
                    dataSource={this.state.data}
                    renderItem={item => (
                        <List.Item>
                            <List.Item.Meta
                                title={item.uname}
                                avatar={
                                    <img
                                        referrerPolicy="no-referrer"
                                        src="https://i2.hdslb.com/bfs/face/9c01434a7af30190393778a1636153e9336474dd.jpg"
                                        alt="face"
                                        style={{width: 100, height: 100}}
                                    />
                                }
                                description={
                                    <Row justify="space-between">
                                        <Col gutter={4} style={{width: '130px'}}>
                                            <Statistic title="直播时长"
                                                       suffix="小时"
                                                       precision={1}
                                                       value={item.duration}
                                            />
                                        </Col>
                                        <Col gutter={4} style={{width: '130px'}}>
                                            <Statistic title="营收"
                                                       suffix="元"
                                                       precision={0}
                                                       value={item.income}
                                            />
                                        </Col>
                                        <Col gutter={4} style={{width: '130px'}}>
                                            <Statistic title="场均付费"
                                                       suffix="人"
                                                       precision={0}
                                                       value={item.avgPaidUser}
                                            />
                                        </Col>
                                        <Col gutter={4} style={{width: '130px'}}>
                                            <Statistic title="场均互动"
                                                       suffix="人"
                                                       precision={0}
                                                       value={item.avgParticipant}
                                            />
                                        </Col>
                                        <Col gutter={4} style={{width: '130px'}}>
                                            <Statistic title="平均同接"
                                                       suffix="人"
                                                       precision={0}
                                                       value={item.avgParticipant}
                                            />
                                        </Col>
                                        <Col gutter={4} style={{width: '130px'}}>
                                            <Statistic title="弹幕总数"
                                                       suffix="条"
                                                       precision={0}
                                                       value={item.danmuCount}
                                            />
                                        </Col>
                                    </Row>
                                }
                            />
                        </List.Item>
                    )}
                />
            </div>
        )
    }
}

export default withRouter(Rank);