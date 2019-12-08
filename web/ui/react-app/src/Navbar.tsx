import React, { FC, useState } from 'react';
import { Link } from '@reach/router';
import {
  Collapse,
  DropdownItem,
  DropdownMenu,
  DropdownToggle,
  Nav,
  Navbar,
  NavbarToggler,
  NavItem,
  NavLink,
  UncontrolledDropdown,
} from 'reactstrap';

const Navigation: FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const toggle = () => setIsOpen(!isOpen);

  return (
    <Navbar className="mb-3" dark color="dark" expand="md" fixed="top">
      <NavbarToggler onClick={toggle} />
      <Link className="pt-0 navbar-brand" to="/ui/playground">
        Prometheus Webhook Dingtalk
      </Link>
      <Collapse isOpen={isOpen} navbar style={{ justifyContent: 'space-between' }}>
        <Nav className="ml-0" navbar>
          <NavItem>
            <NavLink tag={Link} to="/ui/playground">
              Playground
            </NavLink>
          </NavItem>
          <UncontrolledDropdown nav inNavbar>
            <DropdownToggle nav caret>
              Status
            </DropdownToggle>
            <DropdownMenu>
              <DropdownItem tag={Link} to="/ui/status">
                Runtime & Build Information
              </DropdownItem>
              <DropdownItem tag={Link} to="/ui/flags">
                Command-Line Flags
              </DropdownItem>
              <DropdownItem tag={Link} to="/ui/config">
                Configuration
              </DropdownItem>
            </DropdownMenu>
          </UncontrolledDropdown>
        </Nav>
      </Collapse>
    </Navbar>
  );
};

export default Navigation;
