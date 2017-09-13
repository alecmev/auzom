import { PropTypes } from 'react';
import ReactSelect from 'react-select';

import Component from '../utils/Component';

export default class Select extends Component {
  static propTypes = {
    label: PropTypes.string,
    value: PropTypes.any, // eslint-disable-line
    options: PropTypes.array, // eslint-disable-line
    multi: PropTypes.bool,
    onChange: PropTypes.func,
    parentField: PropTypes.string, // stateless
    modified: PropTypes.bool,
  };

  constructor(props) {
    super(props);
    this.state = { value: props.value, inputValue: '' };
  }

  handleChange = (x) => {
    this.setState({ value: x, inputValue: '' });
    this.props.onChange && this.props.onChange(x, this.props.parentField);
  };

  handleInputChange = x => this.setState({ inputValue: x });

  focus() {
    this.r.input.focus();
  }

  render() {
    const { value: stateValue, inputValue } = this.state;
    const fixedStateValue = (
      (Array.isArray(stateValue) && !stateValue.length) ? null : stateValue
    );
    const isEmpty = !fixedStateValue && !inputValue.length;
    return (
      <div className={this.cn({
        d: 'outer',
        s: [isEmpty && 'empty'],
      })}>
        <ReactSelect
          ref={this.rcb('input')}
          placeholder={this.props.label}
          value={this.props.parentField ? this.props.value : fixedStateValue}
          options={this.props.options}
          onChange={this.handleChange}
          onInputChange={this.handleInputChange}
          multi={this.props.multi}
        />
        <div className={this.cn({
          d: 'label',
          s: [isEmpty && 'collapsed'],
        })}>
          {this.props.label}
        </div>
        <div
          className={this.cn({
            d: 'modified',
            s: [!this.props.modified && 'collapsed'],
          })}
          title="modified"
        />
      </div>
    );
  }
}
