import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Modal, Button, Tab } from 'react-bootstrap';
import Axios from 'axios';
import { toast } from 'react-toastify';
import { DragDropContext, Draggable, Droppable } from 'react-beautiful-dnd';

export default class TableSettings extends Component {
    constructor(props) {
        super(props)

        this.state = {
            visible: this.props.data[0],
            hidden: this.props.data[1],
        }
    }

    handleOnDragEnd = (result) => {
        if (!result.destination) return;

        let visible = this.state.visible;
        let hidden = this.state.hidden;

        let from = result.source.droppableId == 'visibleColumns' ? visible.obj.value : hidden.obj.value;
        let to = result.destination.droppableId == 'visibleColumns' ? visible.obj.value : hidden.obj.value;

        let [reorderedItem] = from.splice(result.source.index, 1);
        to.splice(result.destination.index, 0, reorderedItem);

        this.setState({
            visible: visible,
            hidden: hidden
        });
    }

    save = () => {
        var url = 'api/settings/bulk?token=' + window.token;

        Axios.post(url, {
            data: [
                {
                    name: 'visible_columns',
                    value: this.state.visible.obj.value
                },
                {
                    name: 'hidden_columns',
                    value: this.state.hidden.obj.value
                }
            ],
        })
        .then((resp) => {
            toast.success('Table settings updated');
            this.props.refreshConfig();
        })
        .catch((err) => {
            toast.error('Something went wrong');
            console.log(err);
        })
    }

    render() {
        let visible = this.state.visible;
        let hidden = this.state.hidden;

        return (
            <Tab.Content>
                <div>
                    <p>{visible.obj.description}</p>

                    <DragDropContext onDragEnd={this.handleOnDragEnd}>
                        <div className="card pt-4 pb-2 px-4 mb-4">
                            <h4>Visible Columns</h4>
                            <Droppable droppableId="visibleColumns">
                                {(provided) => (
                                    <ul className="visibleColumns pl-0" {...provided.droppableProps} ref={provided.innerRef}>
                                        {visible.obj.value.map((e, i) => {
                                            return (
                                                <Draggable draggableId={e} index={i} key={e}>
                                                    {(provided) => (
                                                        <li className="card bg-secondary py-2 px-3 my-2" key={e} ref={provided.innerRef} {...provided.draggableProps} {...provided.dragHandleProps}>{e}</li>
                                                    )}
                                                </Draggable>
                                            );
                                        })}
                                        {provided.placeholder}
                                    </ul>
                                )}
                            </Droppable>
                        </div>

                        <div className="card pt-4 pb-2 px-4">
                            <h4>Hidden Columns</h4>
                            <Droppable droppableId="hiddenColumns pl-0">
                                {(provided) => (
                                    <ul className="hiddenColumns pl-0" {...provided.droppableProps} ref={provided.innerRef}>
                                        {hidden.obj.value.map((e, i) => {
                                            return (
                                                <Draggable draggableId={e} index={i} key={e}>
                                                    {(provided) => (
                                                        <li className="card bg-secondary py-2 px-3 my-2" key={e} ref={provided.innerRef} {...provided.draggableProps} {...provided.dragHandleProps}>{e}</li>
                                                    )}
                                                </Draggable>
                                            );
                                        })}
                                        {provided.placeholder}
                                    </ul>
                                )}
                            </Droppable>
                        </div>
                    </DragDropContext>

                    <div className="mt-3">
                        <button className="btn btn-primary" onClick={() => { this.save() }}>Save</button>
                    </div>
                </div>
            </Tab.Content>
        );
    }
}

if (document.getElementById('TableSettings')) {
    ReactDOM.render(<TableSettings />, document.getElementById('TableSettings'));
}
