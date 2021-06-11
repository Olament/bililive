import React from 'react';
import 'antd/dist/antd.css';
import './index.css';
import './App.css';
import {Statistic, Row, Col, Modal} from 'antd';
import Broadcast from "./Broadcast";
import VirtualScroller from "virtual-scroller/react";

export default class CardList extends React.Component {
    state = {
        loading: true,
        count: 0,
        data: [],
        isModalVisible: false,
        modalData: {},
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
            <>
                <VirtualScroller
                    id="broadcasts"
                    items={this.state.data}
                    itemComponent={({children}) => (
                        <div
                            style={{padding: '14px 14px 8px'}}
                        >
                            <Broadcast
                                item={children}
                                modalClick={(e, item)=>{
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
                    getColumnsCount={(container) => (Math.min(Math.floor(container.getWidth() / 240), 4))}
                />
                <Modal
                    visible={this.state.isModalVisible}
                    centered={true}
                    onCancel={() => {
                        this.setState({isModalVisible: false})
                    }}
                    footer={null}
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
                        <Col span={10}>
                            <Statistic title="当前人气" value={this.state.modalData.popularity}/>
                        </Col>
                        <Col span={5}>
                            <Statistic title="最高人气" value={this.state.modalData.maxPopularity}/>
                        </Col>
                    </Row>
                    <Row>
                        <Col span={10}>
                            <Statistic title="金瓜子" value={this.state.modalData.goldCoin}/>
                        </Col>
                        <Col span={5}>
                            <Statistic title="银色瓜子" value={this.state.modalData.silverCoin}/>
                        </Col>
                    </Row>
                </Modal>
            </>
        );
    }
}