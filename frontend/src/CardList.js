import React from 'react';
import 'antd/dist/antd.css';
import './index.css';
import './App.css';
import {Statistic, Row, Col, Modal} from 'antd';
import Broadcast from "./Broadcast";
import VirtualScroller from "virtual-scroller/react";
import { LoadingOutlined } from '@ant-design/icons';


export default class CardList extends React.Component {
    state = {
        count: 0,
        data: [],
        isModalVisible: false,
        isLoading: true,
        modalData: {},
    };

    componentDidMount() {
        clearInterval(this.timer)
        this.getData(res => {
            this.setState({
                count: res.count,
                data: res.list,
                isLoading: false,
            })
        })
        this.timer = setInterval(() => {
            this.getData(res => {
                this.setState({
                    count: res.count,
                    data: res.list,
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
            <>
                {this.state.isLoading && <LoadingOutlined style={{ fontSize: 24 }} spin />}
                <VirtualScroller
                    id="broadcasts"
                    items={this.state.data}
                    itemComponent={({children}) => (
                        <div
                            style={{padding: '14px 14px 8px'}}
                        >
                            <Broadcast
                                item={children}
                                modalClick={(e, item) => {
                                    if (e.target.nodeName !== 'A') {
                                        this.setState({
                                            isModalVisible: true,
                                            modalData: item,
                                        })
                                    }
                                }}
                            >
                            </Broadcast>
                        </div>
                    )}
                    getColumnsCount={(container) => (Math.min(Math.floor((container.getWidth() - 50) / 268), 4))}
                />
                <Modal
                    visible={this.state.isModalVisible}
                    centered={true}
                    onCancel={() => {
                        this.setState({isModalVisible: false})
                    }}
                    footer={null}
                    title={this.state.modalData.title}
                >
                    <Row gutter={16}>
                        <Col>
                            <Statistic title="开始时间" value={(() => {
                                let time = new Date(this.state.modalData.livetime)
                                return time.toLocaleString('zh-CN')
                            })()
                            }/>
                        </Col>
                    </Row>
                    <Row>
                        <Col span={12}>
                            <Statistic title="当前人气" value={this.state.modalData.popularity}/>
                        </Col>
                        <Col span={6}>
                            <Statistic title="最高人气" value={this.state.modalData.maxPopularity}/>
                        </Col>
                    </Row>
                    <Row>
                        <Col>
                            <Statistic title="金瓜子" value={this.state.modalData.goldCoin}/>
                        </Col>
                    </Row>
                    <Row>
                        <Col>
                            <Statistic title="银瓜子" value={this.state.modalData.silverCoin}/>
                        </Col>
                    </Row>
                    <Row>
                        <Col span={7}>
                            <Statistic title="付费人数" value={this.state.modalData.goldUser}/>
                        </Col>
                        <Col span={7}>
                            <Statistic title="参与人数" value={this.state.modalData.participant}/>
                        </Col>
                        <Col span={7}>
                            <Statistic title="弹幕数" value={this.state.modalData.danmuCount}/>
                        </Col>
                    </Row>
                </Modal>
            </>
        );
    }
}