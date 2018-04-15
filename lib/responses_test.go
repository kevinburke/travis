package travis

var buildResponse = []byte(`
{
  "@type": "builds",
  "@href": "/repo/kevinburke%2Ftravis/builds?branch.name=master",
  "@representation": "standard",
  "@pagination": {
    "limit": 25,
    "offset": 0,
    "count": 2,
    "is_first": true,
    "is_last": true,
    "next": null,
    "prev": null,
    "first": {
      "@href": "/repo/kevinburke%2Ftravis/builds?branch.name=master",
      "offset": 0,
      "limit": 25
    },
    "last": {
      "@href": "/repo/kevinburke%2Ftravis/builds?branch.name=master",
      "offset": 0,
      "limit": 25
    }
  },
  "builds": [
    {
      "@type": "build",
      "@href": "/build/366635873",
      "@representation": "standard",
      "@permissions": {
        "read": true,
        "cancel": true,
        "restart": true
      },
      "id": 366635873,
      "number": "2",
      "state": "failed",
      "duration": 36,
      "event_type": "push",
      "previous_state": "passed",
      "pull_request_title": null,
      "pull_request_number": null,
      "started_at": "2018-04-14T23:03:34Z",
      "finished_at": "2018-04-14T23:04:10Z",
      "repository": {
        "@type": "repository",
        "@href": "/repo/18699435",
        "@representation": "minimal",
        "id": 18699435,
        "name": "travis",
        "slug": "kevinburke/travis"
      },
      "branch": {
        "@type": "branch",
        "@href": "/repo/18699435/branch/master",
        "@representation": "minimal",
        "name": "master"
      },
      "tag": null,
      "commit": {
        "@type": "commit",
        "@representation": "minimal",
        "id": 109419530,
        "sha": "4208212faa90e5ce7cf7c41bcc17ce091ab56137",
        "ref": "refs/heads/master",
        "message": "initial work on open command",
        "compare_url": "https://github.com/kevinburke/travis/compare/a0890b604823...4208212faa90",
        "committed_at": "2018-04-14T23:03:23Z"
      },
      "jobs": [
        {
          "@type": "job",
          "@href": "/job/366635874",
          "@representation": "minimal",
          "id": 366635874
        }
      ],
      "stages": [

      ],
      "created_by": {
        "@type": "user",
        "@href": "/user/6151",
        "@representation": "minimal",
        "id": 6151,
        "login": "kevinburke"
      },
      "updated_at": "2018-04-14T23:04:10.260Z"
    },
    {
      "@type": "build",
      "@href": "/build/366632359",
      "@representation": "standard",
      "@permissions": {
        "read": true,
        "cancel": true,
        "restart": true
      },
      "id": 366632359,
      "number": "1",
      "state": "passed",
      "duration": 56,
      "event_type": "api",
      "previous_state": null,
      "pull_request_title": null,
      "pull_request_number": null,
      "started_at": "2018-04-14T22:44:31Z",
      "finished_at": "2018-04-14T22:45:27Z",
      "repository": {
        "@type": "repository",
        "@href": "/repo/18699435",
        "@representation": "minimal",
        "id": 18699435,
        "name": "travis",
        "slug": "kevinburke/travis"
      },
      "branch": {
        "@type": "branch",
        "@href": "/repo/18699435/branch/master",
        "@representation": "minimal",
        "name": "master"
      },
      "tag": null,
      "commit": {
        "@type": "commit",
        "@representation": "minimal",
        "id": 109418312,
        "sha": "a0890b60482341a26e549f3c005405d6da0517bc",
        "ref": null,
        "message": "Initial commit",
        "compare_url": "https://github.com/kevinburke/travis/commit/a0890b60",
        "committed_at": "2018-04-14T22:42:14Z"
      },
      "jobs": [
        {
          "@type": "job",
          "@href": "/job/366632360",
          "@representation": "minimal",
          "id": 366632360
        }
      ],
      "stages": [

      ],
      "created_by": {
        "@type": "user",
        "@href": "/user/6151",
        "@representation": "minimal",
        "id": 6151,
        "login": "kevinburke"
      },
      "updated_at": "2018-04-14T22:45:27.750Z"
    }
  ]
}`)

var jobResponse = []byte(`
{
      "@type": "job",
      "@href": "/job/366686566",
      "@representation": "standard",
      "@permissions": {
        "read": true,
        "delete_log": true,
        "cancel": true,
        "restart": true,
        "debug": true
      },
      "id": 366686566,
      "allow_failure": false,
      "number": "9.2",
      "state": "passed",
      "started_at": "2018-04-15T04:42:29Z",
      "finished_at": "2018-04-15T04:43:18Z",
      "build": {
        "@href": "/build/366686564"
      },
      "queue": "builds.ec2",
      "repository": {
        "@type": "repository",
        "@href": "/repo/18699435",
        "@representation": "minimal",
        "id": 18699435,
        "name": "travis",
        "slug": "kevinburke/travis"
      },
      "commit": {
        "@type": "commit",
        "@representation": "minimal",
        "id": 109435658,
        "sha": "15790c98ec36a01b68b19717b576a33b958a7f0e",
        "ref": "refs/heads/wait",
        "message": "Implement rudimentary \"travis wait\" command\n\nget previous build and use better heuristics in shouldPrint",
        "compare_url": "https://github.com/kevinburke/travis/compare/731069ad287a...15790c98ec36",
        "committed_at": "2018-04-15T04:42:17Z"
      },
      "owner": {
        "@type": "user",
        "@href": "/user/6151",
        "@representation": "minimal",
        "id": 6151,
        "login": "kevinburke"
      },
      "stage": null,
      "created_at": "2018-04-15T04:42:26.392Z",
      "updated_at": "2018-04-15T04:43:18.528Z"
    }
`)

var logContent = "travis_fold:start:worker_info\r\u001b[0K\u001b[33;1mWorker information\u001b[0m\nhostname: 35c84bcf-bfa2-4c99-bbe2-127abfbc653b@1.i-0e371a3-production-2-worker-org-ec2.travisci.net\nversion: v3.6.0 https://github.com/travis-ci/worker/tree/170b2a0bb43234479fd1911ba9e4dbcc36dadfad\ninstance: 18ede18 travisci/ci-garnet:packer-1512502276-986baf0 (via amqp)\nstartup: 524.849625ms\ntravis_fold:end:worker_info\r\u001b[0Kmode of ‘/usr/local/clang-5.0.0/bin’ changed from 0777 (rwxrwxrwx) to 0775 (rwxrwxr-x)\r\ntravis_fold:start:system_info\r\u001b[0K\u001b[33;1mBuild system information\u001b[0m\r\nBuild language: go\r\nBuild group: stable\r\nBuild dist: trusty\r\nBuild id: 366686564\r\nJob id: 366686566\r\nRuntime kernel version: 4.14.12-041412-generic\r\ntravis-build version: e0f4abce4\r\n\u001b[34m\u001b[1mBuild image provisioning date and time\u001b[0m\r\nTue Dec  5 20:11:19 UTC 2017\r\n\u001b[34m\u001b[1mOperating System Details\u001b[0m\r\nDistributor ID:\tUbuntu\r\nDescription:\tUbuntu 14.04.5 LTS\r\nRelease:\t14.04\r\nCodename:\ttrusty\r\n\u001b[34m\u001b[1mCookbooks Version\u001b[0m\r\n7c2c6a6 https://github.com/travis-ci/travis-cookbooks/tree/7c2c6a6\r\n\u001b[34m\u001b[1mgit version\u001b[0m\r\ngit version 2.15.1\r\n\u001b[34m\u001b[1mbash version\u001b[0m\r\nGNU bash, version 4.3.11(1)-release (x86_64-pc-linux-gnu)\r\n\u001b[34m\u001b[1mgcc version\u001b[0m\r\ngcc (Ubuntu 4.8.4-2ubuntu1~14.04.3) 4.8.4\r\nCopyright (C) 2013 Free Software Foundation, Inc.\r\nThis is free software; see the source for copying conditions.  There is NO\r\nwarranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.\r\n\r\n\u001b[34m\u001b[1mdocker version\u001b[0m\r\nClient:\r\n Version:      17.09.0-ce\r\n API version:  1.32\r\n Go version:   go1.8.3\r\n Git commit:   afdb6d4\r\n Built:        Tue Sep 26 22:39:28 2017\r\n OS/Arch:      linux/amd64\r\n\u001b[34m\u001b[1mclang version\u001b[0m\r\nclang version 5.0.0 (tags/RELEASE_500/final)\r\nTarget: x86_64-unknown-linux-gnu\r\nThread model: posix\r\nInstalledDir: /usr/local/clang-5.0.0/bin\r\n\u001b[34m\u001b[1mjq version\u001b[0m\r\njq-1.5\r\n\u001b[34m\u001b[1mbats version\u001b[0m\r\nBats 0.4.0\r\n\u001b[34m\u001b[1mshellcheck version\u001b[0m\r\n0.4.6\r\n\u001b[34m\u001b[1mshfmt version\u001b[0m\r\nv2.0.0\r\n\u001b[34m\u001b[1mccache version\u001b[0m\r\nccache version 3.1.9\r\n\r\nCopyright (C) 2002-2007 Andrew Tridgell\r\nCopyright (C) 2009-2011 Joel Rosdahl\r\n\r\nThis program is free software; you can redistribute it and/or modify it under\r\nthe terms of the GNU General Public License as published by the Free Software\r\nFoundation; either version 3 of the License, or (at your option) any later\r\nversion.\r\n\u001b[34m\u001b[1mcmake version\u001b[0m\r\ncmake version 3.9.2\r\n\r\nCMake suite maintained and supported by Kitware (kitware.com/cmake).\r\n\u001b[34m\u001b[1mheroku version\u001b[0m\r\nheroku-cli/6.14.39-addc925 (linux-x64) node-v9.2.0\r\n\u001b[34m\u001b[1mimagemagick version\u001b[0m\r\nVersion: ImageMagick 6.7.7-10 2017-07-31 Q16 http://www.imagemagick.org\r\n\u001b[34m\u001b[1mmd5deep version\u001b[0m\r\n4.2\r\n\u001b[34m\u001b[1mmercurial version\u001b[0m\r\nMercurial Distributed SCM (version 4.2.2)\r\n(see https://mercurial-scm.org for more information)\r\n\r\nCopyright (C) 2005-2017 Matt Mackall and others\r\nThis is free software; see the source for copying conditions. There is NO\r\nwarranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.\r\n\u001b[34m\u001b[1mmysql version\u001b[0m\r\nmysql  Ver 14.14 Distrib 5.6.33, for debian-linux-gnu (x86_64) using  EditLine wrapper\r\n\u001b[34m\u001b[1mopenssl version\u001b[0m\r\nOpenSSL 1.0.1f 6 Jan 2014\r\n\u001b[34m\u001b[1mpacker version\u001b[0m\r\nPacker v1.0.2\r\n\r\nYour version of Packer is out of date! The latest version\r\nis 1.1.2. You can update by downloading from www.packer.io\r\n\u001b[34m\u001b[1mpostgresql client version\u001b[0m\r\npsql (PostgreSQL) 9.6.6\r\n\u001b[34m\u001b[1mragel version\u001b[0m\r\nRagel State Machine Compiler version 6.8 Feb 2013\r\nCopyright (c) 2001-2009 by Adrian Thurston\r\n\u001b[34m\u001b[1msubversion version\u001b[0m\r\nsvn, version 1.8.8 (r1568071)\r\n   compiled Aug 10 2017, 17:20:39 on x86_64-pc-linux-gnu\r\n\r\nCopyright (C) 2013 The Apache Software Foundation.\r\nThis software consists of contributions made by many people;\r\nsee the NOTICE file for more information.\r\nSubversion is open source software, see http://subversion.apache.org/\r\n\r\nThe following repository access (RA) modules are available:\r\n\r\n* ra_svn : Module for accessing a repository using the svn network protocol.\r\n  - with Cyrus SASL authentication\r\n  - handles 'svn' scheme\r\n* ra_local : Module for accessing a repository on local disk.\r\n  - handles 'file' scheme\r\n* ra_serf : Module for accessing a repository via WebDAV protocol using serf.\r\n  - using serf 1.3.3\r\n  - handles 'http' scheme\r\n  - handles 'https' scheme\r\n\r\n\u001b[34m\u001b[1msudo version\u001b[0m\r\nSudo version 1.8.9p5\r\nConfigure options: --prefix=/usr -v --with-all-insults --with-pam --with-fqdn --with-logging=syslog --with-logfac=authpriv --with-env-editor --with-editor=/usr/bin/editor --with-timeout=15 --with-password-timeout=0 --with-passprompt=[sudo] password for %p:  --without-lecture --with-tty-tickets --disable-root-mailer --enable-admin-flag --with-sendmail=/usr/sbin/sendmail --with-timedir=/var/lib/sudo --mandir=/usr/share/man --libexecdir=/usr/lib/sudo --with-sssd --with-sssd-lib=/usr/lib/x86_64-linux-gnu --with-selinux\r\nSudoers policy plugin version 1.8.9p5\r\nSudoers file grammar version 43\r\n\r\nSudoers path: /etc/sudoers\r\nAuthentication methods: 'pam'\r\nSyslog facility if syslog is being used for logging: authpriv\r\nSyslog priority to use when user authenticates successfully: notice\r\nSyslog priority to use when user authenticates unsuccessfully: alert\r\nSend mail if the user is not in sudoers\r\nUse a separate timestamp for each user/tty combo\r\nLecture user the first time they run sudo\r\nRoot may run sudo\r\nAllow some information gathering to give useful error messages\r\nRequire fully-qualified hostnames in the sudoers file\r\nVisudo will honor the EDITOR environment variable\r\nSet the LOGNAME and USER environment variables\r\nLength at which to wrap log file lines (0 for no wrap): 80\r\nAuthentication timestamp timeout: 15.0 minutes\r\nPassword prompt timeout: 0.0 minutes\r\nNumber of tries to enter a password: 3\r\nUmask to use or 0777 to use user's: 022\r\nPath to mail program: /usr/sbin/sendmail\r\nFlags for mail program: -t\r\nAddress to send mail to: root\r\nSubject line for mail messages: *** SECURITY information for %h ***\r\nIncorrect password message: Sorry, try again.\r\nPath to authentication timestamp dir: /var/lib/sudo\r\nDefault password prompt: [sudo] password for %p: \r\nDefault user to run commands as: root\r\nValue to override user's $PATH with: /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin\r\nPath to the editor for use by visudo: /usr/bin/editor\r\nWhen to require a password for 'list' pseudocommand: any\r\nWhen to require a password for 'verify' pseudocommand: all\r\nFile descriptors >= 3 will be closed before executing a command\r\nEnvironment variables to check for sanity:\r\n\tTZ\r\n\tTERM\r\n\tLINGUAS\r\n\tLC_*\r\n\tLANGUAGE\r\n\tLANG\r\n\tCOLORTERM\r\nEnvironment variables to remove:\r\n\tRUBYOPT\r\n\tRUBYLIB\r\n\tPYTHONUSERBASE\r\n\tPYTHONINSPECT\r\n\tPYTHONPATH\r\n\tPYTHONHOME\r\n\tTMPPREFIX\r\n\tZDOTDIR\r\n\tREADNULLCMD\r\n\tNULLCMD\r\n\tFPATH\r\n\tPERL5DB\r\n\tPERL5OPT\r\n\tPERL5LIB\r\n\tPERLLIB\r\n\tPERLIO_DEBUG \r\n\tJAVA_TOOL_OPTIONS\r\n\tSHELLOPTS\r\n\tGLOBIGNORE\r\n\tPS4\r\n\tBASH_ENV\r\n\tENV\r\n\tTERMCAP\r\n\tTERMPATH\r\n\tTERMINFO_DIRS\r\n\tTERMINFO\r\n\t_RLD*\r\n\tLD_*\r\n\tPATH_LOCALE\r\n\tNLSPATH\r\n\tHOSTALIASES\r\n\tRES_OPTIONS\r\n\tLOCALDOMAIN\r\n\tCDPATH\r\n\tIFS\r\nEnvironment variables to preserve:\r\n\tJAVA_HOME\r\n\tTRAVIS\r\n\tCI\r\n\tDEBIAN_FRONTEND\r\n\tXAUTHORIZATION\r\n\tXAUTHORITY\r\n\tPS2\r\n\tPS1\r\n\tPATH\r\n\tLS_COLORS\r\n\tKRB5CCNAME\r\n\tHOSTNAME\r\n\tHOME\r\n\tDISPLAY\r\n\tCOLORS\r\nLocale to use while parsing sudoers: C\r\nDirectory in which to store input/output logs: /var/log/sudo-io\r\nFile in which to store the input/output log: %{seq}\r\nAdd an entry to the utmp/utmpx file when allocating a pty\r\nPAM service name to use\r\nPAM service name to use for login shells\r\nCreate a new PAM session for the command to run in\r\nMaximum I/O log sequence number: 0\r\n\r\nLocal IP address and netmask pairs:\r\n\t172.17.0.2/255.255.0.0\r\n\r\nSudoers I/O plugin version 1.8.9p5\r\n\u001b[34m\u001b[1mgzip version\u001b[0m\r\ngzip 1.6\r\nCopyright (C) 2007, 2010, 2011 Free Software Foundation, Inc.\r\nCopyright (C) 1993 Jean-loup Gailly.\r\nThis is free software.  You may redistribute copies of it under the terms of\r\nthe GNU General Public License <http://www.gnu.org/licenses/gpl.html>.\r\nThere is NO WARRANTY, to the extent permitted by law.\r\n\r\nWritten by Jean-loup Gailly.\r\n\u001b[34m\u001b[1mzip version\u001b[0m\r\nCopyright (c) 1990-2008 Info-ZIP - Type 'zip \"-L\"' for software license.\r\nThis is Zip 3.0 (July 5th 2008), by Info-ZIP.\r\nCurrently maintained by E. Gordon.  Please send bug reports to\r\nthe authors using the web page at www.info-zip.org; see README for details.\r\n\r\nLatest sources and executables are at ftp://ftp.info-zip.org/pub/infozip,\r\nas of above date; see http://www.info-zip.org/ for other sites.\r\n\r\nCompiled with gcc 4.8.2 for Unix (Linux ELF) on Oct 21 2013.\r\n\r\nZip special compilation options:\r\n\tUSE_EF_UT_TIME       (store Universal Time)\r\n\tBZIP2_SUPPORT        (bzip2 library version 1.0.6, 6-Sept-2010)\r\n\t    bzip2 code and library copyright (c) Julian R Seward\r\n\t    (See the bzip2 license for terms of use)\r\n\tSYMLINK_SUPPORT      (symbolic links supported)\r\n\tLARGE_FILE_SUPPORT   (can read and write large files on file system)\r\n\tZIP64_SUPPORT        (use Zip64 to store large files in archives)\r\n\tUNICODE_SUPPORT      (store and read UTF-8 Unicode paths)\r\n\tSTORE_UNIX_UIDs_GIDs (store UID/GID sizes/values using new extra field)\r\n\tUIDGID_NOT_16BIT     (old Unix 16-bit UID/GID extra field not used)\r\n\t[encryption, version 2.91 of 05 Jan 2007] (modified for Zip 3)\r\n\r\nEncryption notice:\r\n\tThe encryption code of this program is not copyrighted and is\r\n\tput in the public domain.  It was originally written in Europe\r\n\tand, to the best of our knowledge, can be freely distributed\r\n\tin both source and object forms from any country, including\r\n\tthe USA under License Exception TSU of the U.S. Export\r\n\tAdministration Regulations (section 740.13(e)) of 6 June 2002.\r\n\r\nZip environment options:\r\n             ZIP:  [none]\r\n          ZIPOPT:  [none]\r\n\u001b[34m\u001b[1mvim version\u001b[0m\r\nVIM - Vi IMproved 7.4 (2013 Aug 10, compiled Nov 24 2016 16:43:18)\r\nIncluded patches: 1-52\r\nExtra patches: 8.0.0056\r\nModified by pkg-vim-maintainers@lists.alioth.debian.org\r\nCompiled by buildd@\r\nHuge version without GUI.  Features included (+) or not (-):\r\n+acl             +farsi           +mouse_netterm   +syntax\r\n+arabic          +file_in_path    +mouse_sgr       +tag_binary\r\n+autocmd         +find_in_path    -mouse_sysmouse  +tag_old_static\r\n-balloon_eval    +float           +mouse_urxvt     -tag_any_white\r\n-browse          +folding         +mouse_xterm     -tcl\r\n++builtin_terms  -footer          +multi_byte      +terminfo\r\n+byte_offset     +fork()          +multi_lang      +termresponse\r\n+cindent         +gettext         -mzscheme        +textobjects\r\n-clientserver    -hangul_input    +netbeans_intg   +title\r\n-clipboard       +iconv           +path_extra      -toolbar\r\n+cmdline_compl   +insert_expand   -perl            +user_commands\r\n+cmdline_hist    +jumplist        +persistent_undo +vertsplit\r\n+cmdline_info    +keymap          +postscript      +virtualedit\r\n+comments        +langmap         +printer         +visual\r\n+conceal         +libcall         +profile         +visualextra\r\n+cryptv          +linebreak       +python          +viminfo\r\n+cscope          +lispindent      -python3         +vreplace\r\n+cursorbind      +listcmds        +quickfix        +wildignore\r\n+cursorshape     +localmap        +reltime         +wildmenu\r\n+dialog_con      -lua             +rightleft       +windows\r\n+diff            +menu            -ruby            +writebackup\r\n+digraphs        +mksession       +scrollbind      -X11\r\n-dnd             +modify_fname    +signs           -xfontset\r\n-ebcdic          +mouse           +smartindent     -xim\r\n+emacs_tags      -mouseshape      -sniff           -xsmp\r\n+eval            +mouse_dec       +startuptime     -xterm_clipboard\r\n+ex_extra        +mouse_gpm       +statusline      -xterm_save\r\n+extra_search    -mouse_jsbterm   -sun_workshop    -xpm\r\n   system vimrc file: \"$VIM/vimrc\"\r\n     user vimrc file: \"$HOME/.vimrc\"\r\n 2nd user vimrc file: \"~/.vim/vimrc\"\r\n      user exrc file: \"$HOME/.exrc\"\r\n  fall-back for $VIM: \"/usr/share/vim\"\r\nCompilation: gcc -c -I. -Iproto -DHAVE_CONFIG_H     -g -O2 -fstack-protector --param=ssp-buffer-size=4 -Wformat -Werror=format-security -U_FORTIFY_SOURCE -D_FORTIFY_SOURCE=1      \r\nLinking: gcc   -Wl,-Bsymbolic-functions -Wl,-z,relro -Wl,--as-needed -o vim        -lm -ltinfo -lnsl  -lselinux  -lacl -lattr -lgpm -ldl    -L/usr/lib/python2.7/config-x86_64-linux-gnu -lpython2.7 -lpthread -ldl -lutil -lm -Xlinker -export-dynamic -Wl,-O1 -Wl,-Bsymbolic-functions      \r\n\u001b[34m\u001b[1miptables version\u001b[0m\r\niptables v1.4.21\r\n\u001b[34m\u001b[1mcurl version\u001b[0m\r\ncurl 7.35.0 (x86_64-pc-linux-gnu) libcurl/7.35.0 OpenSSL/1.0.1f zlib/1.2.8 libidn/1.28 librtmp/2.3\r\n\u001b[34m\u001b[1mwget version\u001b[0m\r\nGNU Wget 1.15 built on linux-gnu.\r\n\u001b[34m\u001b[1mrsync version\u001b[0m\r\nrsync  version 3.1.0  protocol version 31\r\n\u001b[34m\u001b[1mgimme version\u001b[0m\r\nv1.2.0\r\n\u001b[34m\u001b[1mnvm version\u001b[0m\r\n0.33.6\r\n\u001b[34m\u001b[1mperlbrew version\u001b[0m\r\n/home/travis/perl5/perlbrew/bin/perlbrew  - App::perlbrew/0.80\r\n\u001b[34m\u001b[1mphpenv version\u001b[0m\r\nrbenv 1.1.1-25-g6aa70b6\r\n\u001b[34m\u001b[1mrvm version\u001b[0m\r\nrvm 1.29.3 (latest) by Michal Papis, Piotr Kuczynski, Wayne E. Seguin [https://rvm.io]\r\n\u001b[34m\u001b[1mdefault ruby version\u001b[0m\r\nruby 2.4.1p111 (2017-03-22 revision 58053) [x86_64-linux]\r\n\u001b[34m\u001b[1mCouchDB version\u001b[0m\r\ncouchdb 1.6.1\r\n\u001b[34m\u001b[1mElasticSearch version\u001b[0m\r\n5.5.0\r\n\u001b[34m\u001b[1mInstalled Firefox version\u001b[0m\r\nfirefox 56.0.2\r\n\u001b[34m\u001b[1mMongoDB version\u001b[0m\r\nMongoDB 3.4.10\r\n\u001b[34m\u001b[1mPhantomJS version\u001b[0m\r\n2.1.1\r\n\u001b[34m\u001b[1mPre-installed PostgreSQL versions\u001b[0m\r\n9.2.24\r\n9.3.20\r\n9.4.15\r\n9.5.10\r\n9.6.6\r\n\u001b[34m\u001b[1mRabbitMQ Version\u001b[0m\r\n3.6.14\r\n\u001b[34m\u001b[1mRedis version\u001b[0m\r\nredis-server 4.0.6\r\n\u001b[34m\u001b[1mriak version\u001b[0m\r\n2.2.3\r\n\u001b[34m\u001b[1mPre-installed Go versions\u001b[0m\r\n1.7.4\r\n\u001b[34m\u001b[1mant version\u001b[0m\r\nApache Ant(TM) version 1.9.3 compiled on April 8 2014\r\n\u001b[34m\u001b[1mmvn version\u001b[0m\r\nApache Maven 3.5.2 (138edd61fd100ec658bfa2d307c43b76940a5d7d; 2017-10-18T07:58:13Z)\r\nMaven home: /usr/local/maven-3.5.2\r\nJava version: 1.8.0_151, vendor: Oracle Corporation\r\nJava home: /usr/lib/jvm/java-8-oracle/jre\r\nDefault locale: en_US, platform encoding: UTF-8\r\nOS name: \"linux\", version: \"4.4.0-101-generic\", arch: \"amd64\", family: \"unix\"\r\n\u001b[34m\u001b[1mgradle version\u001b[0m\r\n\r\n------------------------------------------------------------\r\nGradle 4.0.1\r\n------------------------------------------------------------\r\n\r\nBuild time:   2017-07-07 14:02:41 UTC\r\nRevision:     38e5dc0f772daecca1d2681885d3d85414eb6826\r\n\r\nGroovy:       2.4.11\r\nAnt:          Apache Ant(TM) version 1.9.6 compiled on June 29 2015\r\nJVM:          1.8.0_151 (Oracle Corporation 25.151-b12)\r\nOS:           Linux 4.4.0-101-generic amd64\r\n\r\n\u001b[34m\u001b[1mlein version\u001b[0m\r\nLeiningen 2.8.1 on Java 1.8.0_151 Java HotSpot(TM) 64-Bit Server VM\r\n\u001b[34m\u001b[1mPre-installed Node.js versions\u001b[0m\r\nv4.8.6\r\nv6.12.0\r\nv6.12.1\r\nv8.9\r\nv8.9.1\r\n\u001b[34m\u001b[1mphpenv versions\u001b[0m\r\n  system\r\n  5.6\r\n* 5.6.32 (set by /home/travis/.phpenv/version)\r\n  7.0\r\n  7.0.25\r\n  7.1\r\n  7.1.11\r\n  hhvm\r\n  hhvm-stable\r\n\u001b[34m\u001b[1mcomposer --version\u001b[0m\r\nComposer version 1.5.2 2017-09-11 16:59:25\r\n\u001b[34m\u001b[1mPre-installed Ruby versions\u001b[0m\r\nruby-2.2.7\r\nruby-2.3.4\r\nruby-2.4.1\r\ntravis_fold:end:system_info\r\u001b[0K\r\nremoved ‘/etc/apt/sources.list.d/basho_riak.list’\r\n\u001b[32;1mNetwork availability confirmed.\u001b[0m\r\n127.0.0.1\tlocalhost\r\n::1\t ip6-localhost ip6-loopback\r\nfe00::0\tip6-localnet\r\nff00::0\tip6-mcastprefix\r\nff02::1\tip6-allnodes\r\nff02::2\tip6-allrouters\r\n172.17.0.4\ttravis-job-kevinburke-travis-366686566.travisci.net travis-job-kevinburke-travis-366686566\r\nW: http://ppa.launchpad.net/couchdb/stable/ubuntu/dists/trusty/Release.gpg: Signature by key 15866BAFD9BCC4F3C1E0DFC7D69548E1C17EAB57 uses weak digest algorithm (SHA1)\r\ntravis_fold:start:git.checkout\r\u001b[0Ktravis_time:start:0f00670c\r\u001b[0K$ git clone --depth=50 --branch=wait https://github.com/kevinburke/travis.git kevinburke/travis\r\nCloning into 'kevinburke/travis'...\r\nremote: Counting objects: 27, done.\u001b[K\r\nremote: Compressing objects:   5% (1/20)   \u001b[K\rremote: Compressing objects:  10% (2/20)   \u001b[K\rremote: Compressing objects:  15% (3/20)   \u001b[K\rremote: Compressing objects:  20% (4/20)   \u001b[K\rremote: Compressing objects:  25% (5/20)   \u001b[K\rremote: Compressing objects:  30% (6/20)   \u001b[K\rremote: Compressing objects:  35% (7/20)   \u001b[K\rremote: Compressing objects:  40% (8/20)   \u001b[K\rremote: Compressing objects:  45% (9/20)   \u001b[K\rremote: Compressing objects:  50% (10/20)   \u001b[K\rremote: Compressing objects:  55% (11/20)   \u001b[K\rremote: Compressing objects:  60% (12/20)   \u001b[K\rremote: Compressing objects:  65% (13/20)   \u001b[K\rremote: Compressing objects:  70% (14/20)   \u001b[K\rremote: Compressing objects:  75% (15/20)   \u001b[K\rremote: Compressing objects:  80% (16/20)   \u001b[K\rremote: Compressing objects:  85% (17/20)   \u001b[K\rremote: Compressing objects:  90% (18/20)   \u001b[K\rremote: Compressing objects:  95% (19/20)   \u001b[K\rremote: Compressing objects: 100% (20/20)   \u001b[K\rremote: Compressing objects: 100% (20/20), done.\u001b[K\r\nremote: Total 27 (delta 8), reused 26 (delta 7), pack-reused 0\u001b[K\r\nUnpacking objects:   3% (1/27)   \rUnpacking objects:   7% (2/27)   \rUnpacking objects:  11% (3/27)   \rUnpacking objects:  14% (4/27)   \rUnpacking objects:  18% (5/27)   \rUnpacking objects:  22% (6/27)   \rUnpacking objects:  25% (7/27)   \rUnpacking objects:  29% (8/27)   \rUnpacking objects:  33% (9/27)   \rUnpacking objects:  37% (10/27)   \rUnpacking objects:  40% (11/27)   \rUnpacking objects:  44% (12/27)   \rUnpacking objects:  48% (13/27)   \rUnpacking objects:  51% (14/27)   \rUnpacking objects:  55% (15/27)   \rUnpacking objects:  59% (16/27)   \rUnpacking objects:  62% (17/27)   \rUnpacking objects:  66% (18/27)   \rUnpacking objects:  70% (19/27)   \rUnpacking objects:  74% (20/27)   \rUnpacking objects:  77% (21/27)   \rUnpacking objects:  81% (22/27)   \rUnpacking objects:  85% (23/27)   \rUnpacking objects:  88% (24/27)   \rUnpacking objects:  92% (25/27)   \rUnpacking objects:  96% (26/27)   \rUnpacking objects: 100% (27/27)   \rUnpacking objects: 100% (27/27), done.\r\n\r\ntravis_time:end:0f00670c:start=1523767361392424653,finish=1523767361683378722,duration=290954069\r\u001b[0K$ cd kevinburke/travis\r\n$ git checkout -qf 15790c98ec36a01b68b19717b576a33b958a7f0e\r\ntravis_fold:end:git.checkout\r\u001b[0K\u001b[33;1mUpdating gimme\u001b[0m\r\ntravis_time:start:2bf83464\r\u001b[0K$ GIMME_OUTPUT=\"$(gimme 1.10 | tee -a $HOME/.bashrc)\" && eval \"$GIMME_OUTPUT\"\r\ngo version go1.10 linux/amd64\r\n\r\ntravis_time:end:2bf83464:start=1523767369754408447,finish=1523767373991027265,duration=4236618818\r\u001b[0K$ export GOPATH=$HOME/gopath\r\n$ export PATH=$HOME/gopath/bin:$PATH\r\n$ mkdir -p $HOME/gopath/src/github.com/kevinburke/travis\r\n$ rsync -az ${TRAVIS_BUILD_DIR}/ $HOME/gopath/src/github.com/kevinburke/travis/\r\n$ export TRAVIS_BUILD_DIR=$HOME/gopath/src/github.com/kevinburke/travis\r\n$ cd $HOME/gopath/src/github.com/kevinburke/travis\r\ntravis_time:start:07bc7177\r\u001b[0K\r\ntravis_time:end:07bc7177:start=1523767374056861250,finish=1523767374063349201,duration=6487951\r\u001b[0Ktravis_fold:start:cache.1\r\u001b[0KSetting up build cache\r\n$ export CASHER_DIR=$HOME/.casher\r\ntravis_time:start:018c5880\r\u001b[0K$ Installing caching utilities\r\n\r\ntravis_time:end:018c5880:start=1523767378721460123,finish=1523767378753868533,duration=32408410\r\u001b[0Ktravis_time:start:068356a4\r\u001b[0K\r\ntravis_time:end:068356a4:start=1523767378760958582,finish=1523767378766083678,duration=5125096\r\u001b[0Ktravis_time:start:1979bb39\r\u001b[0K\u001b[32;1mattempting to download cache archive\u001b[0m\r\n\u001b[32;1mfetching wait/cache-linux-trusty-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855--go-1.10.tgz\u001b[0m\r\n\u001b[32;1mfound cache\u001b[0m\r\n\r\ntravis_time:end:1979bb39:start=1523767378772009241,finish=1523767383206608778,duration=4434599537\r\u001b[0Ktravis_time:start:009ac974\r\u001b[0K\r\ntravis_time:end:009ac974:start=1523767383212869666,finish=1523767383218299734,duration=5430068\r\u001b[0Ktravis_time:start:0afb806a\r\u001b[0K\u001b[32;1madding /home/travis/gopath/pkg to cache\u001b[0m\r\n\u001b[32;1mcreating directory /home/travis/gopath/pkg\u001b[0m\r\n\r\ntravis_time:end:0afb806a:start=1523767383224636226,finish=1523767385412092340,duration=2187456114\r\u001b[0Ktravis_fold:end:cache.1\r\u001b[0K$ gimme version\r\nv1.3.0\r\n$ go version\r\ngo version go1.10 linux/amd64\r\ntravis_fold:start:go.env\r\u001b[0K$ go env\r\nGOARCH=\"amd64\"\r\nGOBIN=\"\"\r\nGOCACHE=\"/home/travis/.cache/go-build\"\r\nGOEXE=\"\"\r\nGOHOSTARCH=\"amd64\"\r\nGOHOSTOS=\"linux\"\r\nGOOS=\"linux\"\r\nGOPATH=\"/home/travis/gopath\"\r\nGORACE=\"\"\r\nGOROOT=\"/home/travis/.gimme/versions/go1.10.linux.amd64\"\r\nGOTMPDIR=\"\"\r\nGOTOOLDIR=\"/home/travis/.gimme/versions/go1.10.linux.amd64/pkg/tool/linux_amd64\"\r\nGCCGO=\"gccgo\"\r\nCC=\"gcc\"\r\nCXX=\"g++\"\r\nCGO_ENABLED=\"1\"\r\nCGO_CFLAGS=\"-g -O2\"\r\nCGO_CPPFLAGS=\"\"\r\nCGO_CXXFLAGS=\"-g -O2\"\r\nCGO_FFLAGS=\"-g -O2\"\r\nCGO_LDFLAGS=\"-g -O2\"\r\nPKG_CONFIG=\"pkg-config\"\r\nGOGCCFLAGS=\"-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build183711775=/tmp/go-build -gno-record-gcc-switches\"\r\ntravis_fold:end:go.env\r\u001b[0KUsing Go 1.5 Vendoring, not checking for Godeps\r\ntravis_fold:start:install\r\u001b[0Ktravis_time:start:0d25134b\r\u001b[0K$ true\r\n\r\ntravis_time:end:0d25134b:start=1523767385481690929,finish=1523767385486656964,duration=4966035\r\u001b[0Ktravis_fold:end:install\r\u001b[0Ktravis_fold:start:before_script\r\u001b[0Ktravis_time:start:0f7b04b0\r\u001b[0K$ go get ./...\r\n\r\ntravis_time:end:0f7b04b0:start=1523767385492851385,finish=1523767393768352575,duration=8275501190\r\u001b[0Ktravis_fold:end:before_script\r\u001b[0Ktravis_time:start:034329c0\r\u001b[0K$ make race-test\r\ngo test -race ./...\r\nok  \tgithub.com/kevinburke/travis\t1.012s\r\nok  \tgithub.com/kevinburke/travis/lib\t1.014s\r\n\r\ntravis_time:end:034329c0:start=1523767393774481507,finish=1523767397085103955,duration=3310622448\r\u001b[0K\r\n\u001b[32;1mThe command \"make race-test\" exited with 0.\u001b[0m\r\ntravis_fold:start:cache.2\r\u001b[0Kstore build cache\r\ntravis_time:start:1ef8c6cc\r\u001b[0K\r\ntravis_time:end:1ef8c6cc:start=1523767397091762665,finish=1523767397096930176,duration=5167511\r\u001b[0Ktravis_time:start:04d26eca\r\u001b[0K\u001b[32;1mnothing changed, not updating cache\u001b[0m\r\n\r\ntravis_time:end:04d26eca:start=1523767397102928123,finish=1523767398297508622,duration=1194580499\r\u001b[0Ktravis_fold:end:cache.2\r\u001b[0K\r\nDone. Your build exited with 0.\r\n"
