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

import type {JSX} from 'react';
import {Box, Card as MuiCard, Typography, styled, Button} from '@wso2/oxygen-ui';
import {SignUp} from '@asgardeo/react';
import {useNavigate, useSearchParams} from 'react-router';

const Card = styled(MuiCard)(({theme}) => ({
  display: 'flex',
  flexDirection: 'column',
  alignSelf: 'center',
  width: '100%',
  padding: theme.spacing(4),
  gap: theme.spacing(2),
  backdropFilter: 'blur(10px)',
  WebkitBackdropFilter: 'blur(10px)',
  border: '1px solid rgba(255, 255, 255, 0.1) !important',
  background: 'rgba(215, 215, 215, 0.04)',
  boxShadow:
    '0 5px 10px 0 rgba(6, 6, 14, 0.1), 0 0 0 0 rgba(199, 211, 234, 0.01) inset, 0 0 0 0 rgba(199, 211, 234, 0.12) inset',
  [theme.breakpoints.up('sm')]: {
    width: '450px',
  },
  ...theme.applyStyles('dark', {
    boxShadow: 'hsla(220, 30%, 5%, 0.5) 0px 5px 15px 0px, hsla(220, 25%, 10%, 0.08) 0px 15px 35px -5px',
  }),
}));

export type SignUpBoxProps = Record<string, unknown>;

export default function SignUpBox(): JSX.Element {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  return (
    <Card variant="outlined">
      <Box sx={{display: 'flex', flexDirection: 'column', width: '100%', gap: 2}}>
        <Typography component="h1" variant="h5" sx={{width: '100%', mb: 2}}>
          Create Account
        </Typography>

        <SignUp />

        <Typography sx={{textAlign: 'center', mt: 3}}>
          Already have an account?{' '}
          <Button
            variant="text"
            onClick={() => {
              const currentParams = searchParams.toString();
              const signInUrl = currentParams ? `/sign-in?${currentParams}` : '/sign-in';
              // eslint-disable-next-line @typescript-eslint/no-floating-promises
              navigate(signInUrl);
            }}
            sx={{
              p: 0,
              minWidth: 'auto',
              textTransform: 'none',
              color: 'primary.main',
              textDecoration: 'underline',
              '&:hover': {
                textDecoration: 'underline',
                backgroundColor: 'transparent',
              },
            }}
          >
            Sign in
          </Button>
        </Typography>
      </Box>
    </Card>
  );
}
