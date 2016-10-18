import moment from 'moment';
import { PropTypes } from 'react';
import ReactDatetime from 'react-datetime';

import Component from '../utils/Component';

export default class Datetime extends Component {
  static propTypes = {
    label: PropTypes.string,
    value: PropTypes.string,
    justDate: PropTypes.bool,
    onChange: PropTypes.func,
    parentField: PropTypes.string, // stateless
    modified: PropTypes.bool,
    open: PropTypes.bool,
  };

  constructor(props) {
    super(props);
    !this.props.parentField && (this.state = { value: props.value });
  }

  handleChange = (x) => {
    let y;
    if (x) {
      y = x.toISOString();
    }

    !this.props.parentField && this.setState({ value: y });
    this.props.onChange && this.props.onChange(y, this.props.parentField);
  };

  handleClear = () => this.handleChange(null);

  focus() {
    this.r.datetime.openCalendar();
  }

  render() {
    const value = this.props.parentField ? this.props.value : this.state.value;
    return (
      <div className={this.cni()}>
        <ReactDatetime
          ref={this.rcb('datetime')}
          className={this.cn({
            d: 'input',
            s: [!value && 'empty'],
          })}
          inputProps={{
            placeholder: this.props.label,
            readOnly: true,
          }}
          dateFormat={this.props.justDate ? 'D MMMM YYYY' : "D MMM 'YY"}
          timeFormat={this.props.justDate ? false : '[at] HH:mm [your time]'}
          value={value && moment(value)}
          onChange={this.handleChange}
          open={this.props.open}
          closeOnSelect
        />
        <div className={this.cn({
          d: 'label',
          s: [!value && 'collapsed'],
        })}>
          {this.props.label}
        </div>
        <div className={this.cn({
          d: 'modified',
          s: [!this.props.modified && 'collapsed'],
        })} title="modified" />
        <button className={this.cn({
          d: 'clear',
          s: [!value && 'hidden'],
        })} title="clear" onClick={this.handleClear}>×</button>
      </div>
    );
  }
}
