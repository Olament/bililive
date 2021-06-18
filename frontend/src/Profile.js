import React from 'react';
import {List, Statistic, Row, Col, Typography, Card, PageHeader, Tag} from 'antd';
import 'antd/dist/antd.css';
import './App.css'
import {withRouter} from "react-router";


const {Title, Paragraph} = Typography;

class Profile extends React.Component {
    state = {
        data: [],
        loading: true,
    }

    componentDidMount() {
        this.getData(this.props.match.params.uid,
            (res) => (this.setState({
                data: res,
                loading: false
            })))
    }

    getData = (uid, callback) => {
        fetch('https://livevup.com/api/broadcast?uid=' + uid)
            .then(res => res.json())
            .then(res => callback(res))
    }

    timeToDuration = (start, end) => {
        const diff = (new Date(end)) - (new Date(start));
        return (diff / 36e5).toFixed(1);
    }

    average = (list) => (
        list.reduce((total, item) => (total + item), 0) / Math.max(list.length, 1)
    )

    render() {
        console.log(this.state.data)
        return (
            <PageHeader
                className="site-page-header"
                onBack={() => (this.props.history.goBack())}
                title={
                    <a
                        href={"https://space.bilibili.com/" + this.props.match.params.uid}
                        target="_blank"
                        rel="noreferrer"
                    >
                        {this.state.data.length > 0 ? this.state.data[0].uname : "未知"}
                    </a>
                }
            >
                <Card>
                    <Row gutter={36} justify="space-around">
                        <Col gutter={12}>
                            <Statistic
                                title="人均付费"
                                value={
                                    this.average(this.state.data.map(
                                        (item)=>((item.goldCoin / 1000) / Math.max(1, item.goldUser))))
                                }
                                precision={2}
                                suffix="元"
                                style={{width: '120px'}}
                            />
                        </Col>
                        <Col gutter={12}>
                            <Statistic
                                title="人均弹幕"
                                value={
                                    this.average(this.state.data.map(
                                        (item)=>(item.danmuCount / Math.max(1, item.participant))))
                                }
                                suffix="条"
                                precision={2}
                                style={{width: '120px'}}
                            />
                        </Col>
                        <Col gutter={12}>
                            <Statistic
                                title="平均同接"
                                value={
                                    this.average(this.state.data.reduce((acc, item)=>
                                        (acc.concat(item.participantTrend)), []))
                                }
                                suffix="人"
                                precision={0}
                                style={{width: '120px'}}
                            />
                        </Col>
                    </Row>
                </Card>
                <List
                    itemLayout="horizontal"
                    size="large"
                    header={<Title level={4}>近十场直播数据</Title>}
                    loading={this.state.loading}
                    dataSource={this.state.data}
                    renderItem={item => (
                        <List.Item>
                            <List.Item.Meta
                                title={
                                    <div>
                                        <Title level={4}>{item.title}</Title>
                                        <Paragraph type="secondary">
                                            {(new Date(item.livetime)).toLocaleString('zh-CN')}
                                        </Paragraph>
                                    </div>
                                }
                                description={
                                    <Row gutter={36} justify="space-between">
                                        <Col gutter={6} style={{width: '130px'}}>
                                            <Statistic title="直播时长"
                                                       suffix="小时"
                                                       value={this.timeToDuration(item.livetime, item.endTime)}
                                            />
                                        </Col>
                                        <Col gutter={6} style={{width: '130px'}}>
                                            <Statistic title="营收" suffix="元" value={Math.floor(item.goldCoin / 1000)}/>
                                        </Col>
                                        <Col gutter={6} style={{width: '130px'}}>
                                            <Statistic title="付费人数" value={item.goldUser}/>
                                        </Col>
                                        <Col gutter={6} style={{width: '130px'}}>
                                            <Statistic title="互动人数" value={item.participant}/>
                                        </Col>
                                        <Col gutter={6} style={{width: '130px'}}>
                                            <Statistic title="弹幕总数" value={item.danmuCount}/>
                                        </Col>
                                        <Col gutter={6} style={{width: '130px'}}>
                                            <Statistic title="人气峰值" value={item.maxPopularity}/>
                                        </Col>
                                    </Row>
                                }
                            />
                        </List.Item>
                    )}
                />
            </PageHeader>
        )
    }
}

export default withRouter(Profile);