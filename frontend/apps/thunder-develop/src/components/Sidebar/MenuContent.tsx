/**
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import {NavLink} from 'react-router';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Stack from '@mui/material/Stack';
import {Blocks, House, Info, LayoutGrid, MessageCircleQuestionMark, Settings, User, UsersRound} from 'lucide-react';
import useNavigation from '@/layouts/contexts/useNavigation';

const mainListItems = [
  {
    id: 'home',
    text: 'Home',
    icon: <House size={16} />,
    category: 'Dashboard',
    path: '/',
  },
  {
    id: 'users',
    text: 'Users',
    icon: <UsersRound size={16} />,
    category: 'Dashboard',
    path: '/users',
  },
  {
    id: 'user-types',
    text: 'User Types',
    icon: <User size={16} />,
    category: 'Dashboard',
    path: '/user-types',
  },
  {
    id: 'analytics',
    text: 'Integrations',
    icon: <Blocks size={16} />,
    category: 'Dashboard',
    path: '/integrations',
  },
  {
    id: 'tasks',
    text: 'Applications',
    icon: <LayoutGrid size={16} />,
    category: 'Dashboard',
    path: '/applications',
  },
];

const secondaryListItems = [
  {
    id: 'settings',
    text: 'Settings',
    icon: <Settings size={14} />,
    category: 'Settings',
    path: '/settings',
  },
];

export default function MenuContent() {
  const {currentPage, setCurrentPage} = useNavigation();

  const handleListItemClick = (item: {id: string; text: string; category: string}) => {
    setCurrentPage({
      id: item.id,
      text: item.text,
      category: item.category,
    });
  };

  return (
    <Stack sx={{flexGrow: 1, p: 1, justifyContent: 'space-between'}}>
      <List dense>
        {mainListItems.map((item) => (
          <ListItem key={item.id} disablePadding sx={{display: 'block'}}>
            <ListItemButton
              component={NavLink}
              to={item.path}
              selected={currentPage.id === item.id}
              onClick={() => handleListItemClick(item)}
            >
              <ListItemIcon>{item.icon}</ListItemIcon>
              <ListItemText primary={item.text} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
      <List dense>
        {secondaryListItems.map((item) => (
          <ListItem key={item.id} disablePadding sx={{display: 'block'}}>
            <ListItemButton
              component={NavLink}
              to={item.path}
              selected={currentPage.id === item.id}
              onClick={() => handleListItemClick(item)}
            >
              <ListItemIcon>{item.icon}</ListItemIcon>
              <ListItemText primary={item.text} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Stack>
  );
}
